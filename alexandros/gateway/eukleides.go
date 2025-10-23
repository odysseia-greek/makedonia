package gateway

import (
	"context"

	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/makedonia/alexandros/graph/model"
	pbe "github.com/odysseia-greek/makedonia/eukleides/proto"
)

func (a *AlexandrosHandler) pushToEukleides(update *pbe.CountCreationRequest) {

	var collector pbe.CountCreationRequestSet
	collector.Request = append(collector.Request, update)

	err := a.CounterStreamer.Send(&collector)
	if err != nil {
		logging.Error(err.Error())
	}
}

func (a *AlexandrosHandler) TopFive(ctx context.Context) (*model.EukleidesTopFiveResponse, error) {

	in := &pbe.TopFiveRequest{}
	response, err := a.Counter.RetrieveTopFive(ctx, in)

	if err != nil {
		return nil, err
	}

	var out model.EukleidesTopFiveResponse

	for _, resp := range response.TopFive {
		topFive := model.EukleidesTopFive{
			ServiceName: resp.ServiceName,
			Word:        resp.Word,
			LastUsed:    &resp.LastUsed,
			Count:       int32(resp.Count),
		}

		out.TopFive = append(out.TopFive, &topFive)
	}
	return &out, nil
}
