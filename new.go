package lib

func _new(options []Option) (*Claims, error) {
	claims := &Claims{}
	for _, opt := range options {
		err := opt(claims)
		if err != nil {
			return nil, err
		}
	}
	return claims, nil
}

func New(options ...Option) (*Claims, error) {
	claims, err := _new(options)
	if err != nil {
		return nil, err
	}

	claims.Calculate()

	return claims, nil
}
