package handlers

import (
	"fmt"
	"log"
	"matterpoll-bot/internal/entities"
	"matterpoll-bot/internal/services"
	"net/http"
	"strings"

	"github.com/mattermost/mattermost-server/v6/model"
)

// CreatePoll обрабатывает HTTP-запрос и разбирает полученные параметры в соответствии с примером:
// "text": строка в формате `/poll-create "Question" "Option1" "Option2" ...`,
// где Question — вопрос для голосвания, а Option1, Option2  — варианты для голоса.
// Обработчик разбирает параметр "text", чтобы извлечь вопрос и варианты ответа.
// Если создание голосования прошло успешно, возвращается сообщение с результатом
// Если формат параметра "text" некорректен, возвращается сообщение об ошибке с примером правильного формата.
// Создает новые опросы, закрепляя за ними id создателя.
func CreatePoll(s *services.PollService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		text := r.Form.Get("text")
		args := strings.Split(text, `" "`)
		if len(args) < 2 {
			w.Write([]byte("**Invalid format!** *Example*: `/poll-create \"Question\" \"Option1\" \"Option2\" ...`"))
			return
		}

		question := strings.Trim(args[0], `"`)
		options := args[1:]

		voices := make(map[string]int32, len(options))
		for i := range options {
			options[i] = strings.Trim(options[i], `"`)
			voices[options[i]] = 0
		}

		optionsStr := strings.Join(options, "` `")
		optionsStr = "`" + optionsStr + "`"

		id := model.NewId()
		userId := r.Form.Get("user_id")
		if userId == "" {
			http.Error(w, "'user_id' is empty in the form data", http.StatusBadRequest)
			return
		}

		poll := &entities.Poll{
			PollId:   id,
			Question: question,
			Options:  voices,
			Creator:  userId,
			Voters:   map[string]bool{},
			Closed:   false,
		}

		err := s.CreatePoll(poll)
		if err != nil {
			if userErr, ok := err.(*entities.UserError); ok {
				w.Write([]byte(userErr.Error()))
				return
			}
			log.Println(err)
			http.Error(w, "failed to create Poll", http.StatusInternalServerError)
			return
		}

		channelId := r.Form.Get("channel_id")
		if channelId == "" {
			http.Error(w, "'channel_id' is empty in the form data", http.StatusBadRequest)
			return
		}

		post := &model.Post{ChannelId: channelId, Message: fmt.Sprintf("**Poll created!** *Poll_ID*: `%s` *Question*: `%s` *Options*: %s", id, question, optionsStr)}
		_, resp, err := s.Bot.CreatePost(post)
		if err != nil {
			if userErr, ok := err.(*entities.UserError); ok {
				w.Write([]byte(userErr.Error()))
				return
			}

			log.Println(err)
			http.Error(w, "failed to create Poll", http.StatusInternalServerError)
		}

		if resp == nil || resp.StatusCode != 201 {
			w.Write([]byte(fmt.Sprintf("failed to get team: unexpected status code %d", resp.StatusCode)))
		}
	}
}

// Vote обрабатывает HTTP-запрос для голосования в опросе.
// Ожидается, что запрос будет содержать параметры формы:
// "text": строка в формате `"Poll_ID" "Option"`, где Poll_ID — идентификатор опроса, а Option — выбранный вариант.
// Обработчик разбирает параметр "text", чтобы извлечь идентификатор опроса и вариант ответа для голоса.
// Если голосование успешно, возвращается сообщение с результатом.
// Если формат параметра "text" некорректен, возвращается сообщение об ошибке с примером правильного формата.
// В случае ошибки возвращается соответствующее сообщение об ошибке или статус HTTP 500 для внутренних ошибок сервера.
func Vote(s *services.PollService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		text := r.Form.Get("text")
		args := strings.Split(text, `" "`)
		if len(args) != 2 {
			w.Write([]byte("**Invalid format!** *Example*: `/poll-vote \"Poll_ID\" \"Option\"`"))
			return
		}

		pollId := strings.Trim(args[0], `"`)
		option := strings.Trim(args[1], `"`)

		userId := r.Form.Get("user_id")
		if userId == "" {
			http.Error(w, "'user_id' is empty in the form data", http.StatusBadRequest)
			return
		}

		voice := &entities.Voice{PollId: pollId, UserId: userId, Option: option}
		msg, err := s.Vote(voice)
		if err != nil {
			if userErr, ok := err.(*entities.UserError); ok {
				w.Write([]byte(userErr.Error()))
				return
			}
			log.Println(err)
			http.Error(w, "failed to vote", http.StatusInternalServerError)
		}

		w.Write([]byte(msg))
	}
}

// GetPollResults обрабатывает HTTP-запрос для получения результатов опроса.
// Ожидается, что запрос будет содержать параметр формы:
// "text": строка в формате `"Poll_ID"`, где Poll_ID — идентификатор опроса.
// Обработчик разбирает параметр "text", чтобы извлечь идентификатор опроса.
// Если формат параметра "text" некорректен, возвращается сообщение об ошибке с примером правильного формата.
// Если получение результатов прошло успешно, возвращается сообщение с результатами опроса.
// В случае ошибки возвращается соответствующее сообщение об ошибке или статус HTTP 500 для внутренних ошибок сервера.
func GetPollResults(s *services.PollService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		text := r.Form.Get("text")
		args := strings.Split(text, `" "`)
		if len(args) != 1 {
			w.Write([]byte("**Invalid format!** *Example*: `/poll-results \"Poll_ID\"`"))
			return
		}

		pollId := strings.Trim(args[0], `"`)
		msg, err := s.GetPollResult(pollId)
		if err != nil {
			if userErr, ok := err.(*entities.UserError); ok {
				w.Write([]byte(userErr.Error()))
				return
			}

			log.Println(err)
			http.Error(w, "failed to get poll results", http.StatusInternalServerError)
			return
		}

		w.Write([]byte(msg))
	}
}

// ClosePoll обрабатывает HTTP-запрос для закрытия опроса.
// Ожидается, что запрос будет содержать следующие параметры формы:
// "text": строка в формате `"Poll_ID"`, где Poll_ID — идентификатор опроса.
// Обработчик разбирает параметр "text", чтобы извлечь идентификатор опроса.
// Если формат параметра "text" некорректен, возвращается сообщение об ошибке с примером правильного формата.
// Если получение результатов прошло успешно, возвращается сообщение с результатами опроса.
// В случае ошибки возвращается соответствующее сообщение об ошибке или статус HTTP 500 для внутренних ошибок сервера.
func ClosePoll(s *services.PollService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		text := r.Form.Get("text")
		args := strings.Split(text, `" "`)
		if len(args) != 1 {
			w.Write([]byte("**Неверный формат!** *Пример*: `/poll-close \"Poll_ID\"`"))
			return
		}

		pollId := strings.Trim(args[0], `"`)
		userId := r.Form.Get("user_id")
		if userId == "" {
			http.Error(w, "'user_id' is empty in the form data", http.StatusBadRequest)
			return
		}

		msg, err := s.ClosePoll(pollId, userId)
		if err != nil {
			if userErr, ok := err.(*entities.UserError); ok {
				w.Write([]byte(userErr.Error()))
				return
			}

			log.Println(err)
			http.Error(w, "failed to close poll", http.StatusInternalServerError)
			return
		}

		w.Write([]byte(msg))
	}
}

// DeletePoll обрабатывает HTTP-запрос для удаления опроса.
// Этот обработчик ожидает, что запрос будет содержать следующие параметры формы:
// Ожидается, что запрос будет содержать следующие параметры формы:
// "text": строка в формате `"Poll_ID"`, где Poll_ID — идентификатор опроса.
// Обработчик разбирает параметр "text", чтобы извлечь идентификатор опроса.
// Если формат параметра "text" некорректен, возвращается сообщение об ошибке с примером правильного формата.
// Если получение результатов прошло успешно, возвращается сообщение с результатами опроса.
// В случае ошибки возвращается соответствующее сообщение об ошибке или статус HTTP 500 для внутренних ошибок сервера.
func DeletePoll(s *services.PollService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		text := r.Form.Get("text")
		args := strings.Split(text, `" "`)
		if len(args) != 1 {
			w.Write([]byte("**Invalid format!** *Example*: `/poll-close \"Poll_ID\"`"))
			return
		}

		pollId := strings.Trim(args[0], `"`)
		userId := r.Form.Get("user_id")
		if userId == "" {
			http.Error(w, "'user_id' is empty in the form data", http.StatusBadRequest)
			return
		}

		msg, err := s.DeletePoll(pollId, userId)
		if err != nil {
			if userErr, ok := err.(*entities.UserError); ok {
				w.Write([]byte(userErr.Error()))
				return
			}

			log.Println(err)
			http.Error(w, "failed to delete", http.StatusInternalServerError)
			return
		}

		w.Write([]byte(msg))
	}
}
