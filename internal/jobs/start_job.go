package jobs

func StartJob() error {
	jobCtx := NewJob()

	if err := jobCtx.RegJob("0 */1 * * * *", NewDemoJobBuilder().Build()); err != nil {
		return err
	}

	go jobCtx.Serve()

	return nil
}
