package handler

import (
	"context"
	"time"

	// "github.com/micro/dev/model"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/model"

	"github.com/gosimple/slug"

	posts "posts/proto"
	// tags "tags/proto"
)

type Posts struct {
	db model.Model
	// Tags tags.TagsService
}

func NewPosts() *Posts {
	createdIndex := model.ByEquality("created")
	createdIndex.Order.Type = model.OrderTypeDesc

	// create model
	db := model.NewModel(
		model.WithIndexes(model.ByEquality("slug"), createdIndex),
	)
	// register the post instance
	db.Register(new(posts.Post))

	return &Posts{
		db: db,
		// Tags: tagsService,
	}
}

func (p *Posts) Save(ctx context.Context, req *posts.SaveRequest, rsp *posts.SaveResponse) error {
	logger.Info("Received Posts.Save request")
	// validate
	if len(req.Id) == 0 {
		return errors.BadRequest("posts.save.MissingId", "missing id")
	}

	// read post
	recs := []*posts.Post{}
	q := model.QueryEquals("id", req.Id)
	q.Order.Type = model.OrderTypeUnordered
	if err := p.db.Read(q, &recs); err != nil {
		return errors.InternalServerError("posts.save.Read", "failed to read post with id")
	}

	// if no existing post is found, create one
	if len(recs) == 0 {
		if len(req.Title) == 0 {
			return errors.BadRequest("posts.save.MissingTitle", "missing title")
		}
		pslug := slug.Make(req.Title)
		// get post with slug
		recsWithSlug := []*posts.Post{}
		if err := p.db.Read(model.QueryEquals("slug", pslug), &recsWithSlug); err != nil {
			return errors.InternalServerError("posts.save.Read", "failed to read post with slug")
		}
		if len(recsWithSlug) > 0 {
			return errors.BadRequest("posts.save.SlugCheck", "existing slug")
		}

		post := &posts.Post{
			Id:       req.Id,
			Title:    req.Title,
			Content:  req.Content,
			Tags:     req.Tags,
			Slug:     pslug,
			Created:  req.Timestamp,
			Metadata: req.Metadata,
			Image:    req.Image,
		}
		if post.Created == 0 {
			post.Created = time.Now().Unix()
		}
		if err := p.db.Create(post); err != nil {
			return errors.InternalServerError("posts.save.Create", "failed to create post")
		}
		rsp.Id = post.Id
		return nil
	}

	// existing post
	post := recs[0]
	post.Updated = req.Timestamp
	post.Metadata = req.Metadata
	post.Image = req.Image
	if post.Created == 0 {
		post.Created = time.Now().Unix()
	}
	if len(req.Title) > 0 {
		post.Title = req.Title
		post.Slug = slug.Make(req.Title)
	}
	if len(req.Slug) > 0 {
		post.Slug = req.Slug
	}
	if len(req.Tags) > 0 {
		if req.Tags[0] == "" {
			post.Tags = []string{}
		} else {
			post.Tags = req.Tags
		}
	}

	// get post with slug
	recsWithSlug := []*posts.Post{}
	if err := p.db.Read(model.QueryEquals("slug", post.Slug), &recsWithSlug); err != nil {
		return errors.InternalServerError("posts.save.Read", "failed to read post with slug")
	}
	if len(recsWithSlug) > 0 && post.Id != recsWithSlug[0].Id {
		return errors.BadRequest("posts.save.SlugCheck", "existing slug")
	}

	// update
	if err := p.db.Create(post); err != nil {
		return errors.InternalServerError("posts.save.Create", "failed to update post")
	}
	rsp.Id = post.Id

	return nil
}

// Index returns the posts index without content
func (p *Posts) Index(ctx context.Context, req *posts.IndexRequest, rsp *posts.IndexResponse) error {
	logger.Info("Received Posts.Index request")
	q := model.QueryEquals("created", nil)
	q.Order.Type = model.OrderTypeDesc
	q.Offset = req.Offset
	q.Limit = req.Limit

	return p.db.Read(q, &rsp.Posts)

	// var recs []*posts.Post

	// // read
	// if err := p.db.Read(q, &recs); err != nil {
	// 	return errors.InternalServerError("posts.index.Read", "failed to read posts")
	// }

	// for _, post := range recs {
	// 	rsp.Posts = append(rsp.Posts, post)
	// }

	// return nil
}

// Query currently only supports read by slug or timestamp, no listing.
func (p *Posts) Query(ctx context.Context, req *posts.QueryRequest, rsp *posts.QueryResponse) error {
	logger.Info("Received Posts.Query request")
	var q model.Query
	if len(req.Slug) > 0 {
		logger.Infof("Reading by slug: %v", req.Slug)
		q = model.QueryEquals("slug", req.Slug)
	} else if len(req.Id) > 0 {
		logger.Infof("Reading by id: %v", req.Id)
		q = model.QueryEquals("id", req.Id)
		q.Order.Type = model.OrderTypeUnordered
	} else {
		q = model.QueryEquals("created", nil)
		q.Order.Type = model.OrderTypeDesc
		q.Limit = 20
		if req.Limit > 0 {
			q.Limit = req.Limit
		}
		q.Offset = req.Offset
		logger.Infof("Reading offset: %v, limit: %v", q.Limit, q.Offset)
	}

	return p.db.Read(q, &rsp.Posts)
}

func (p *Posts) Delete(ctx context.Context, req *posts.DeleteRequest, rsp *posts.DeleteResponse) error {
	logger.Info("Received Posts.Delete request")
	// validate
	if len(req.Id) == 0 {
		return errors.BadRequest("posts.delete.MissingId", "missing id")
	}
	q := model.QueryEquals("id", req.Id)
	q.Order.Type = model.OrderTypeUnordered
	return p.db.Delete(q)
}
