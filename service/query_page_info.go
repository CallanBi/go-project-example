package service

// service 层需要做的事情：
// 1. 参数校验：对于外部传入的参数，需要做校验，比如topicId是否大于0
// 2. 业务逻辑：对于业务逻辑，需要做一些处理，比如查询topic信息，查询post列表，组装pageInfo
// 3. 错误处理：对于错误，需要做一些处理，比如参数校验失败，查询topic信息失败，查询post列表失败，返回错误信息

import (
	"errors"
	"sync"

	"github.com/Moonlight-Zhao/go-project-example/repository"
)

type PageInfo struct {
	Topic    *repository.Topic
	PostList []*repository.Post
}

func QueryPageInfo(topicId int64) (*PageInfo, error) {
	return NewQueryPageInfoFlow(topicId).Do()
}

func NewQueryPageInfoFlow(topId int64) *QueryPageInfoFlow {
	return &QueryPageInfoFlow{
		topicId: topId,
	}
}

type QueryPageInfoFlow struct {
	topicId  int64
	pageInfo *PageInfo

	topic   *repository.Topic
	posts   []*repository.Post
}

func (f *QueryPageInfoFlow) Do() (*PageInfo, error) {
	if err := f.checkParam(); err != nil {
		return nil, err
	}
	if err := f.prepareInfo(); err != nil {
		return nil, err
	}
	if err := f.packPageInfo(); err != nil {
		return nil, err
	}
	return f.pageInfo, nil
}

func (f *QueryPageInfoFlow) checkParam() error {
	if f.topicId <= 0 {
		return errors.New("topic id must be larger than 0")
	}
	return nil
}

func (f *QueryPageInfoFlow) prepareInfo() error {
	//获取topic信息
	// 由于获取topic信息和获取post列表是两个独立的操作，前置依赖都是TopicId，所以可以考虑使用goroutine并行执行
	var wg sync.WaitGroup
	// 两个goroutine，所以需要wg.Add(2)
	wg.Add(2)
	go func() {
		defer wg.Done()
		topic := repository.NewTopicDaoInstance().QueryTopicById(f.topicId)
		f.topic = topic
	}()
	//获取post列表
	go func() {
		defer wg.Done()
		posts := repository.NewPostDaoInstance().QueryPostsByParentId(f.topicId)
		f.posts = posts
	}()
	wg.Wait()
	return nil
}

func (f *QueryPageInfoFlow) packPageInfo() error {
	f.pageInfo = &PageInfo{
		Topic:    f.topic,
		PostList: f.posts,
	}
	return nil
}
