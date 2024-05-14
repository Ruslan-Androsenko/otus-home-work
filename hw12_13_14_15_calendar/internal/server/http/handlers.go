package internalhttp

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/storage"
)

var ctx = context.Background()

func (s *Server) createEventHandler(w http.ResponseWriter, r *http.Request) {
	params, err := loadParams[EventParams](r)
	if !errors.Is(err, io.EOF) {
		logg.Errorf("%v", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	event, errParse := params.getEvent()
	if event == nil || errParse != nil {
		logg.Errorf("Failed to parse create params: %v, Error: %v", params, errParse)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if err = s.app.CreateEvent(ctx, *event); err != nil {
		logg.Errorf("Failed to create new event: %v, Error: %v", *event, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := DataResponse{"state": "success", "data": "Ok"}
	w.WriteHeader(http.StatusCreated)
	SendResponse(w, response, createEventMethod)
}

func (s *Server) updateEventHandler(w http.ResponseWriter, r *http.Request) {
	params, err := loadParams[EventParams](r)
	if !errors.Is(err, io.EOF) {
		logg.Errorf("%v", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	event, errParse := params.getEvent()
	if event == nil || errParse != nil {
		logg.Errorf("Failed to parse update params: %v, Error: %v", params, errParse)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if err = s.app.UpdateEvent(ctx, params.ID, *event); err != nil {
		logg.Errorf("Failed to update eventID: %s, event: %v, Error: %v", params.ID, *event, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := DataResponse{"state": "success", "data": "Ok"}
	SendResponse(w, response, updateEventMethod)
}

func (s *Server) deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	params, err := loadParams[EventParams](r)
	if !errors.Is(err, io.EOF) {
		logg.Errorf("%v", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err = s.app.DeleteEvent(ctx, params.ID); err != nil {
		logg.Errorf("Failed to delete eventID: %v, Error: %v", params.ID, err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	response := DataResponse{"state": "success", "data": "Ok"}
	SendResponse(w, response, deleteEventMethod)
}

func (s *Server) getEventHandler(w http.ResponseWriter, r *http.Request) {
	params, err := loadParams[EventParams](r)
	if !errors.Is(err, io.EOF) {
		logg.Errorf("%v", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	event, err := s.app.GetEventByID(params.ID)
	if err != nil {
		errIs := err.Error()
		if errIs == storage.ErrSQLNoRows {
			response := DataResponse{"state": "not found", "data": nil}
			w.WriteHeader(http.StatusNotFound)
			SendResponse(w, response, getEventMethod)
			return
		} else {
			logg.Errorf("Failed to get event by id: %v, Error: %v", params.ID, err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	response := DataResponse{"state": "success", "data": event}
	SendResponse(w, response, getEventMethod)
}

func (s *Server) getEventsOfDayHandler(w http.ResponseWriter, r *http.Request) {
	getEventsListByDate(w, r, s.app.GetEventsOfDay, getEventOfDayMethod)
}

func (s *Server) getEventsOfWeekHandler(w http.ResponseWriter, r *http.Request) {
	getEventsListByDate(w, r, s.app.GetEventsOfWeek, getEventOfWeekMethod)
}

func (s *Server) getEventsOfMonthHandler(w http.ResponseWriter, r *http.Request) {
	getEventsListByDate(w, r, s.app.GetEventsOfMonth, getEventOfMonthMethod)
}

func getEventsListByDate(
	w http.ResponseWriter,
	r *http.Request,
	getEventsFn func(time time.Time) ([]storage.Event, error),
	caption string,
) {
	params, err := loadParams[EventParams](r)
	if !errors.Is(err, io.EOF) {
		logg.Errorf("%v", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	eventDate, errParse := params.getDate()
	if errParse != nil {
		logg.Errorf("Failed to parse get events params: %v, Error: %v", params, errParse)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	events, err := getEventsFn(eventDate)
	if err != nil {
		logg.Errorf("Failed to get events of date: %v, Error: %v", eventDate, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := DataResponse{"state": "success", "data": events}

	if events == nil {
		response = DataResponse{"state": "not found", "data": nil}
		w.WriteHeader(http.StatusNotFound)
	}

	SendResponse(w, response, caption)
}
