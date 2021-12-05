package sumoapp

//func CompileApp(basePath string) ([]byte, error) {
//	app := NewApplication()
//
//	//Load each app stream and all objects within the stream's app
//	streams, err := LoadAppStreams(basePath)
//	if err != nil {
//		return nil, err
//	}
//
//	//For each app stream, merge with the previous one. The order here is important!
//	for _, stream := range streams {
//		app.Merge(stream.Application)
//	}
//
//	if err := app.Build(); err != nil {
//		return nil, err
//	}
//
//	//Compile into JSON return
//	jsonByteString, err := json.Marshal(app)
//	if err != nil {
//		return nil, err
//	}
//
//	return jsonByteString, nil
//}
