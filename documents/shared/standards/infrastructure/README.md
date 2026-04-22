# README 

coding-standard.mdはこのディレクトリで書きコマンドを打ってで生成すること

```
docker run --rm -v ".:/work:ro" -w /work --entrypoint sh mikefarah/yq:4 convert.sh > coding-standard.md
```
