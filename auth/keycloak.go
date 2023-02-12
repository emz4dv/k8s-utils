package auth

import (
	"context"
	"k8s-utils/conf"
	"log"

	"github.com/Nerzal/gocloak/v10"
)


func GetToken(conf *conf.Config) string {
	client := conf.Auth
	clientGoCloak := gocloak.NewClient(client.URL)
	ctx := context.Background()
	token, err := clientGoCloak.Login(ctx, client.ClientID, client.ClientSecret, client.Realm, client.User, client.Password)
	if err != nil {
		log.Fatal("error:", err)
	} 

	return token.AccessToken
}	
