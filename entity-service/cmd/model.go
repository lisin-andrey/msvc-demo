package cmd

import (
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/lisin-andrey/msvc-demo/common/pkg/tools"
	"github.com/lisin-andrey/msvc-demo/entity-service/pkg/config"
	"github.com/lisin-andrey/msvc-demo/entity-service/pkg/repository"
)

// IRepInstanceHandler - for safe getting IEntityRepository
type IRepInstanceHandler interface {
	GetRepInstance() repository.IEntityRepository
}

// RepInstanceHandler - implementation of IRepInstanceHandler
type RepInstanceHandler struct {
	Config *config.Config
	// Use it by call setupSignalHandler
	muRepEntityInit     sync.Mutex
	_singletonRepEntity repository.IEntityRepository
	_repEntityInitFlag  uint32
}

// GetRepInstance - implement IRepInstanceHandler interface
// With health check support
func (h *RepInstanceHandler) GetRepInstance() repository.IEntityRepository {
	if atomic.LoadUint32(&h._repEntityInitFlag) == 1 {
		return h._singletonRepEntity
	}

	h.muRepEntityInit.Lock()
	defer h.muRepEntityInit.Unlock()

	if h._repEntityInitFlag == 0 {
		var err error
		// Create entity repository
		h._singletonRepEntity, err = repository.NewEntityRepositoryByConfig(h.Config)
		if err != nil {
			tools.Errorfln("Can't create EntityRepository. Error: [%s]", err.Error())
			return nil
		}
		atomic.StoreUint32(&h._repEntityInitFlag, 1)
		tools.Debugln("EntityRepository was created")
	}
	return h._singletonRepEntity
}

// RestHandleFunc - Used format of http HandleFuncs in the project
type RestHandleFunc func(repository.IEntityRepository, http.ResponseWriter, *http.Request)

// RestHandler - wrapper of http.Handler
type RestHandler struct {
	RepHandler IRepInstanceHandler
	RestHandleFunc
}

// implement http.Handler interface
func (h *RestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rep := h.RepHandler.GetRepInstance()
	if rep != nil {
		h.RestHandleFunc(rep, w, r)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}
