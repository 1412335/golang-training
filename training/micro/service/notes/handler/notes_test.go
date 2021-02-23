package handler

import (
	"context"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/errors"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	notes "notes/proto"
	"testing"
)

func createClientService() (notes.NotesService, error) {
	// create and initialise a new service
	srv := service.New()

	// create the proto client
	client := notes.NewNotesService("notes", srv.Client())

	return client, nil
}

func TestNotes_List(t *testing.T) {
	client, _ := createClientService()
	require.NotNil(t, client)
	{
		brsp, err := client.List(context.Background(), &notes.ListRequest{})
		require.NoError(t, err)
		require.NotNil(t, brsp)

		reqs := []notes.CreateRequest{
			{
				Title: "Test",
				Text:  "Test",
			},
			{
				Title: "Test 1",
				Text:  "Test",
			},
		}
		for _, req := range reqs {
			rsp, err := client.Create(context.Background(), &req)
			require.NoError(t, err)
			require.NotNil(t, rsp)
		}
		rsp, err := client.List(context.Background(), &notes.ListRequest{})
		require.NoError(t, err)
		require.NotNil(t, rsp)
		require.Len(t, rsp.Notes, len(reqs)+len(brsp.Notes))
	}
}

func TestNotes_Create(t *testing.T) {
	client, _ := createClientService()
	require.NotNil(t, client)
	{
		req := &notes.CreateRequest{
			Title: "Test",
			Text:  "Test",
		}
		rsp, err := client.Create(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, rsp)
		require.NotNil(t, rsp.Id)
	}
	{
		_, err := client.Create(context.Background(), nil)
		require.Error(t, err)
		e := errors.Parse(err.Error())
		require.Equal(t, int32(http.StatusInternalServerError), e.Code)
	}
	{
		_, err := client.Create(context.Background(), &notes.CreateRequest{})
		require.Error(t, err)
		e := errors.Parse(err.Error())
		require.Equal(t, int32(http.StatusBadRequest), e.Code)
	}
}

func TestNotes_Update(t *testing.T) {

}

func TestNotes_UpdateStream(t *testing.T) {
	client, _ := createClientService()
	require.NotNil(t, client)
	{
		brsp, err := client.List(context.Background(), &notes.ListRequest{})
		require.NoError(t, err)
		require.NotNil(t, brsp)

		// create stream
		stream, err := client.UpdateStream(context.Background())
		require.NoError(t, err)
		require.NotNil(t, stream)

		// loop & send into update stream
		for _, note := range brsp.Notes {
			e := stream.Send(&notes.UpdateRequest{
				Id:    note.Id,
				Title: note.Title + "(update)",
				Text:  note.Text + "(update)",
			})
			require.NoError(t, e)
		}
		// close stream & receive resp
		rsp, err := stream.CloseAndRecv()
		require.Equal(t, io.EOF, err)
		require.NotNil(t, rsp)
	}
}

func TestNotes_Delete(t *testing.T) {

}
