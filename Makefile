include $(GOROOT)/src/Make.$(GOARCH)

TARG=db
GOFILES=db.go classic.go util.go result.go set.go connection.go statement.go doc.go

include $(GOROOT)/src/Make.pkg
