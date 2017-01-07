// Package webfinger is a server implementation of the webfinger specification. This
// is a general-case package which provides the HTTP handlers and interfaces
// for adding webfinger support for your system and resources.
//
// The simplest way to use this is to call webfinger.Default() and
// then register the object as an HTTP handler:
//
//		myResolver = ...
// 		wf := webfinger.Default(myResolver{})
//		wf.NotFoundHandler = // the rest of your app
//		http.ListenAndService(":8080", wf)
//
// However, you can also register the specific webfinger handler to a path. This should
// work on any router that supports net/http.
//
//		myResolver = ...
// 		wf := webfinger.Default(myResolver{})
//		http.Handle(webfinger.WebFingerPath, http.HandlerFunc(wf.Webfinger))
//		http.ListenAndService(":8080", nil)
//
// In either case, the handlers attached to the webfinger service get invoked as
// needed.
package webfinger
