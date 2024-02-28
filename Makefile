default:
	go build -o grepcidr-reverse .
debian: default
	mkdir -p ./debian/grepcidr-reverse/usr/local/bin
	go build -o ./debian/grepcidr-reverse/usr/local/bin/grepcidr-reverse .
	cd debian; dpkg-deb --build grepcidr-reverse
install: default
	cp grepcidr-reverse /usr/local/bin/grepcidr-reverse
clean:
	rm -f grepcidr-reverse
	rm -fr debian/grepcidr-reverse/usr
	rm -f grepcidr-reverse.deb
