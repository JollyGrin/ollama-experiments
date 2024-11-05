provide a filepath and it will load it into a conversational agent. Ask via line number

`go build -o askcode`
`sudo mv askcode /usr/local/bin/`
`sudo chmod +x /usr/local/bin/askcode`

### usage

(relative to current directory)
`askcode -- filename.go`

will open a prompt

```
You:
```

Enter your question and it will respond:

```
Model: 
```

Give it a few seconds and it will begin to respond
