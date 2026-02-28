# 중요 순서

- ctx 처리 더 잘할 것

- `fmt` 를 `log`으로 변환하기

- service 구조체가 값이 빠져도 gen이 가능하게

- 패키지 이름 바꾸기 (codegen)

- go.mod 추가 여부 설정 (+ 레포 생성 여부)

- README 다듬기
# 고민

- `build` 와 `generate`를 이미 cli로 제공하고 있는데 굳이 설정을 `task`로 할 이유가 있는지 고민하기

# Future

- bulk update (xlsx) : 외부 패키지 필요 + json에 비해 어려움

- command 기능 고도화하기 : 기능이 더 필요해지면 할 것

- 탬플릿에서 ast로 변환하기 : 보안 관련, 급선무가 아님

- 더 나은 명명 (en-name) + 카테고리화 : 초반에 추가했으나, 중요한 API는 이름이 잘 구분되어있고, en 생성 퀄리티가 낮아서 나중에 다시 하기로

- `generate` 진행과정 표시

- 태그 config 추가