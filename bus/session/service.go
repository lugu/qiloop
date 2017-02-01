package session

import (
	"fmt"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/net"
)

type passiveService struct {
	wrapper  bus.ServiceWrapper
	endpoint net.EndPoint
	handler  int
	service  uint32
	object   uint32
}

func (w *passiveService) filter() net.Filter {
	return func(hdr *net.Header) (matched bool, keep bool) {
		if hdr.Service == w.service && hdr.Object == w.object {
			if _, ok := w.wrapper[hdr.Action]; ok {
				return true, true
			}
		}
		return false, true
	}
}

func (w *passiveService) send(typ uint8, id, action uint32, data []byte) error {
	h := net.NewHeader(typ, w.service, w.object, action, id)
	m := net.NewMessage(h, data)
	return w.endpoint.Send(m)
}

func (w *passiveService) consumer() net.Consumer {
	return func(msg *net.Message) error {
		if h, ok := w.wrapper[msg.Header.Action]; ok {
			var err2 error // returns the first error encountered
			response, err := h(w, msg.Payload)
			if response != nil {
				err2 = w.send(net.Reply, msg.Header.ID, msg.Header.Action, response)
			}
			if err != nil {
				return err
			}
			return err2
		}
		return fmt.Errorf("action %d not wrapped", msg.Header.Action)
	}
}

func (w *passiveService) addHandler() {
	w.handler = w.endpoint.AddHandler(w.filter(), w.consumer())
}

func (w *passiveService) removeHandler() error {
	if w.handler == 0 {
		return fmt.Errorf("Service %d handler already removed", w.service)
	}
	if err := w.endpoint.RemoveHandler(w.handler); err != nil {
		return fmt.Errorf("cannot remove handler for service %d: %s", w.service, err)
	}
	w.handler = 0
	return nil
}

func (w *passiveService) Emit(actionID uint32, data []byte) error {
	// FIXME: what shall the id be ?
	return w.send(net.Event, 0, actionID, data)
}

func (w *passiveService) Unregister() error {
	return w.removeHandler()
}

func (w *passiveService) ServiceID() uint32 {
	return w.service
}

func (w *passiveService) ObjectID() uint32 {
	return w.object
}

func NewService(service, object uint32, endpoint net.EndPoint,
	wrapper bus.ServiceWrapper) bus.Service {
	s := &passiveService{
		endpoint: endpoint,
		service:  service,
		object:   object,
		wrapper:  wrapper,
	}
	s.addHandler()
	return s
}
