package bot

import tele "gopkg.in/telebot.v3"

type heap map[string]tele.HandlerFunc

func (h heap) Handle(endpoint string, hf tele.HandlerFunc) {
	h[endpoint] = hf
}

func (h heap) Get(endpoint string) tele.HandlerFunc {
	if val, ok := h[endpoint]; ok {
		return val
	}

	return nil
}
