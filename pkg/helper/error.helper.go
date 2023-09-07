package helper

func ErrorHelperPanic(err error) {
	if err != nil {
		panic(err)
		return
	}
}
