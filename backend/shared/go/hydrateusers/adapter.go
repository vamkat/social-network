package userhydrate

// type userHydrateAdapter struct {
// 	u *models.User
// }

// func (a userHydrateAdapter) GetUserId() int64 {
// 	return a.u.UserId.Int64()
// }

// func (a userHydrateAdapter) SetUser(user models.User) {
// 	*(a.u) = user
// }

// func (h *UserHydrator) HydrateUserSlice(ctx context.Context, users []models.User) error {
// 	adapters := make([]models.HasUser, len(users))
// 	for i := range users {
// 		adapters[i] = userHydrateAdapter{u: &users[i]}
// 	}
// 	return h.HydrateUsers(ctx, adapters)
// }
