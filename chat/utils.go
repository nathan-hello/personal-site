package chat

import (
	"net/http"

	"github.com/nathan-hello/personal-site/auth"
)

func GetClientState(r *http.Request) auth.ClientState {
	claims, ok := r.Context().Value(ClaimsContextKey).(*CustomClaims)

	if claims == nil {
		state := DefaultClientState
		// leaving room in the future to add things to state even if claims is nil
		return state
	}

	return ClientState{IsAuthed: ok, Username: claims.Username, UserId: claims.ID}
}
