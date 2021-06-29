package mongodb

import "go.mongodb.org/mongo-driver/mongo"

func IsUniqueIndexViolation(err error) bool {
	if merr, ok := err.(mongo.WriteException); ok {
		if merr.HasErrorCode(11000) { // unique index violation
			return true
		}
	}
	return false
}
