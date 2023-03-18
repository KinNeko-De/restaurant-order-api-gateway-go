package document

import (
	grpccontext "context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	apidocument "github.com/kinneko-de/test-api-contract/golang/kinnekode/document/grpc/v1"
	apiprotobuf "github.com/kinneko-de/test-api-contract/golang/kinnekode/protobuf"
	"io"
	"reflect"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

const Deadline = time.Duration(6000) * time.Millisecond

const ParamDocumentId = "documentId"

func GetDocumentById(context *gin.Context) {
	documentId, err := ReadRequestedDocumentId(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c, dialError := grpc.Dial("localhost:5649", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if dialError != nil {
		context.AbortWithError(http.StatusServiceUnavailable, dialError)
		return
	}
	client := apidocument.NewDocumentServiceClient(c)

	requestDocumentId, _ := apiprotobuf.ToProtobuf(documentId)
	request := apidocument.DownloadDocumentRequest{
		DocumentId: requestDocumentId,
	}

	callContext, cancelFunc := grpccontext.WithDeadline(grpccontext.Background(), GetDeadline())
	defer cancelFunc()

	stream, clientErr := client.DownloadDocument(callContext, &request)
	if clientErr != nil {
		context.AbortWithError(http.StatusServiceUnavailable, clientErr)
		return
	}

	firstRequest, streamErr := stream.Recv()
	if streamErr != nil {
		context.AbortWithError(http.StatusServiceUnavailable, streamErr)
		return
	}

	_, ok := firstRequest.File.(*apidocument.DownloadDocumentResponse_Metadata)
	if !ok {
		context.AbortWithError(http.StatusInternalServerError, errors.New("FileCase of type 'apidocument.DownloadDocumentResponse_Metadata' expected. Actual value is "+reflect.TypeOf(firstRequest.File).String()+"."))
		return
	}
	var metadata = firstRequest.GetMetadata()
	context.Header("Content-Type", metadata.MediaType)
	context.Header("Content-Length", strconv.FormatUint(metadata.Size, 10))

	for {
		current, done, requestErr := ReadNextRequest(stream)
		if done {
			return
		}
		if requestErr != nil {
			context.AbortWithError(http.StatusServiceUnavailable, requestErr)
			return
		}
		if SomethingElseThanChunkWasSent(current) {
			context.AbortWithError(http.StatusInternalServerError, errors.New("FileCase of type 'apidocument.DownloadDocumentResponse_Chunk' expected. Actual value is "+reflect.TypeOf(current.File).String()+"."))
			return
		}

		var chunk = current.GetChunk()
		_, bodyWriteErr := context.Writer.Write(chunk)
		if bodyWriteErr != nil {
			context.AbortWithError(http.StatusInternalServerError, bodyWriteErr)
			return
		}
	}
}

func ReadRequestedDocumentId(context *gin.Context) (uuid.UUID, error) {
	paramId := context.Param(ParamDocumentId)
	documentId, err := uuid.Parse(paramId)
	if err != nil {
		return uuid.UUID{}, err
	}
	return documentId, nil
}

func SomethingElseThanChunkWasSent(current *apidocument.DownloadDocumentResponse) bool {
	_, ok := current.File.(*apidocument.DownloadDocumentResponse_Chunk)
	if !ok {
		return true
	}
	return false
}

func ReadNextRequest(stream apidocument.DocumentService_DownloadDocumentClient) (*apidocument.DownloadDocumentResponse, bool, error) {
	current, err := stream.Recv()

	if err == io.EOF {
		err := stream.CloseSend()
		if err != nil {
			return nil, true, err
		}
		return nil, true, nil
	}

	if err != nil {
		return nil, false, err
	}
	return current, false, nil
}

func GetDeadline() time.Time {
	return time.Now().Add(Deadline)
}
