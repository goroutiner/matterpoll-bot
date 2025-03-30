package handlers

import (
	"fmt"
	"log"
	"matterpoll-bot/entities"
	"matterpoll-bot/internal/services"
	"net/http"
	"strings"

	"github.com/mattermost/mattermost-server/v6/model"
)

func CreatePoll(s *services.PollService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		text := r.Form.Get("text")
		args := strings.Split(text, `" "`)
		if len(args) < 2 {
			w.Write([]byte("**Invalid format!** *Example*: `/poll-create \"Question\" \"Option1\" \"Option2\" ... \"OptionN\"`"))
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
			http.Error(w, "Failed to create Poll", http.StatusInternalServerError)
			return
		}

		channelId := r.Form.Get("channel_id")
		post := &model.Post{ChannelId: channelId, Message: fmt.Sprintf("**Poll created!** *Poll_ID*: `%s` *Question*: `%s` *Options*: %s", id, question, optionsStr)}
		if _, _, err = entities.Bot.CreatePost(post); err != nil {
			if userErr, ok := err.(*entities.UserError); ok {
				w.Write([]byte(userErr.Error()))
				return
			}
			log.Println(err)
			http.Error(w, "Failed to create Poll", http.StatusInternalServerError)
			return
		}
	}
}

func Vote(s *services.PollService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		text := r.Form.Get("text")
		args := strings.Split(text, `" "`)
		if len(args) != 2 {
			w.Write([]byte("**Invalid format!** *Example*: `/poll-vote \"Poll_ID\" \"Option\"`"))
			return
		}

		pollId := strings.Trim(args[0], `"`)
		userId := r.Form.Get("user_id")
		option := strings.Trim(args[1], `"`)

		voice := &entities.Voice{PollId: pollId, UserId: userId, Option: option}
		msg, err := s.Vote(voice)
		if err != nil {
			if userErr, ok := err.(*entities.UserError); ok {
				w.Write([]byte(userErr.Error()))
				return
			}
			log.Println(err)
			http.Error(w, "Failed to vote", http.StatusInternalServerError)
		}

		w.Write([]byte(msg))
	}
}

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
			http.Error(w, "Failed to vote", http.StatusInternalServerError)
			return
		}

		w.Write([]byte(msg))
	}
}

func ClosePoll(s *services.PollService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		text := r.Form.Get("text")
		args := strings.Split(text, `" "`)
		if len(args) != 1 {
			w.Write([]byte("**Invalid format!** *Example*: `/poll-close \"Poll_ID\"`"))
			return
		}

		pollId := strings.Trim(args[0], `"`)
		userId := r.Form.Get("user_id")
		msg, err := s.ClosePoll(pollId, userId)
		if err != nil {
			if userErr, ok := err.(*entities.UserError); ok {
				w.Write([]byte(userErr.Error()))
				return
			}

			log.Println(err)
			http.Error(w, "Failed to close poll", http.StatusInternalServerError)
			return
		}

		w.Write([]byte(msg))
	}
}

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
		msg, err := s.DeletePoll(pollId, userId)
		if err != nil {
			if userErr, ok := err.(*entities.UserError); ok {
				w.Write([]byte(userErr.Error()))
				return
			}

			log.Println(err)
			http.Error(w, "Failed to delete", http.StatusInternalServerError)
			return
		}

		w.Write([]byte(msg))
	}
}
