package handle

import (
	"fmt"

	"github.com/cloudedcat/finance-bot/bot"
	"github.com/cloudedcat/finance-bot/calculator"
	"github.com/cloudedcat/finance-bot/log"
	"github.com/cloudedcat/finance-bot/model"
	tb "gopkg.in/tucnak/telebot.v2"
)

type calculateHandler struct {
	calc   calculator.Service
	logger log.Logger
}

func (h *calculateHandler) handle(bot bot.Bot, m *tb.Message) {
	logInfo := formLogInfo(m, "Calculate")
	groupID := model.GroupID(m.Chat.ID)
	finalDebts, err := h.calc.CalculateDebts(groupID)
	if err != nil {
		bot.SendInternalError(m.Chat, logInfo)
		h.logger.IfErrorw(err, "failed to calculate debts", logInfo...)
		return
	}
	bot.Send(m.Chat, h.formMessage(finalDebts), logInfo)
}

func (h *calculateHandler) formMessage(debts []calculator.FinalDebt) (resp string) {
	if len(debts) == 0 {
		resp = "there ain't debts"
	}
	resp = "list of debts:\n"
	for _, debt := range debts {
		resp += fmt.Sprintf("	@%s -> @%s - %.2f",
			debt.Borrower.Alias, debt.Lender.Alias, debt.Amount)
	}
	return resp
}

// Calculate shows debt for each borrower
func Calculate(bot bot.Bot, calc calculator.Service, logger log.Logger) {
	handler := &calculateHandler{
		calc:   calc,
		logger: logger,
	}
	bot.Handle("/calc", notPrivateOnlyMiddleware(handler.handle))
}
