package handler

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"

	notes "notes/proto"
)

const storePrefix = "notes/"

type Notes struct{}

func (n *Notes) List(ctx context.Context, req *notes.ListRequest, rsp *notes.ListResponse) error {
	recs, err := store.Read(storePrefix, store.Prefix(storePrefix))
	if err != nil {
		logger.Errorf("Error reading records from store: %v", err)
		return errors.InternalServerError("notes.List.Unknown", "error reading records from store")
	}
	rsp.Notes = make([]*notes.Note, len(recs))
	for i, record := range recs {
		var note notes.Note
		if err := json.Unmarshal(record.Value, &note); err != nil {
			logger.Errorf("Error unmarshaling note: %v", err)
			return errors.InternalServerError("notes.List.Unknown", "error unmarshaling note")
		}
		rsp.Notes[i] = &note
	}
	return nil
}

func (n *Notes) Create(ctx context.Context, req *notes.CreateRequest, rsp *notes.CreateResponse) error {
	logger.Info("Received Notes.Create request")
	// validate request
	if len(req.Title) == 0 {
		return errors.BadRequest("notes.Create.MissingTitle", "missing title")
	}

	// note
	note := &notes.Note{
		Id:      uuid.New().String(),
		Created: time.Now().Unix(),
		Title:   req.Title,
		Text:    req.Text,
	}

	// marshal note to bytes
	bytes, err := json.Marshal(note)
	if err != nil {
		logger.Errorf("error marshalling note: %v", err)
		return errors.InternalServerError("notes.Create.Unknown", "error marshalling note")
	}

	// write to store
	key := storePrefix + note.Id
	if err := store.Write(&store.Record{Key: key, Value: bytes}); err != nil {
		logger.Errorf("error writing to store: %v", err)
		return errors.InternalServerError("notes.Create.Unknown", "error writing to store")
	}

	// return
	rsp.Id = note.Id
	return nil
}

func (n *Notes) Update(ctx context.Context, req *notes.UpdateRequest, rsp *notes.UpdateResponse) error {
	// validate request
	if len(req.Id) == 0 {
		return errors.BadRequest("notes.Update.MissingId", "missing id")
	}
	if len(req.Title) == 0 {
		return errors.BadRequest("notes.Update.MissingTitle", "missing title")
	}

	// read note from store
	key := storePrefix + req.Id
	recs, err := store.Read(key)
	if err == store.ErrNotFound {
		return errors.BadRequest("notes.Update.InvalidId", "not found note id")
	} else if err != nil {
		logger.Errorf("error reading notes: %v", err)
		return errors.InternalServerError("notes.Update.Unknown", "error reading note")
	}

	var note notes.Note
	if err := json.Unmarshal(recs[0].Value, &note); err != nil {
		logger.Errorf("error unmarshalling note: %v", err)
		return errors.InternalServerError("notes.Update.Unknown", "error unmarshalling note")
	}
	// update
	note.Title = req.Title
	note.Text = req.Text

	// marshal note to bytes
	bytes, err := json.Marshal(note)
	if err != nil {
		logger.Errorf("error marshalling note: %v", err)
		return errors.InternalServerError("notes.Update.Unknown", "error marshalling note")
	}

	// write to store
	if err := store.Write(&store.Record{Key: key, Value: bytes}); err != nil {
		logger.Errorf("error writing to store: %v", err)
		return errors.InternalServerError("notes.Update.Unknown", "error writing to store")
	}

	// return
	rsp.Id = note.Id
	return nil
}

func (n *Notes) UpdateStream(ctx context.Context, stream notes.Notes_UpdateStreamStream) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			logger.Errorf("error reading from stream: %v", err)
			return errors.InternalServerError("notes.UpdateStream.Unknown", "error reading from stream")
		}
		if err := n.Update(ctx, req, &notes.UpdateResponse{}); err != nil {
			return err
		}
	}
}

func (n *Notes) Delete(ctx context.Context, req *notes.DeleteRequest, rsp *notes.DeleteResponse) error {
	// validate request
	if len(req.Id) == 0 {
		return errors.BadRequest("notes.Delete.MissingId", "missing id")
	}
	// delete note
	if err := store.Delete(storePrefix + req.Id); err == store.ErrNotFound {
		return errors.NotFound("notes.Delete.InvalidId", "not found note id")
	} else if err != nil {
		logger.Errorf("error deleting note: %v", err)
		return errors.InternalServerError("notes.Delete.Unknown", "error deleting note")
	}
	// return
	rsp.Id = req.Id
	return nil
}
