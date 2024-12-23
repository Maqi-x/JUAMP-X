.PHONY: all clean

all:
	python3 scripts/build.py

clean:
	rm -rf bin/