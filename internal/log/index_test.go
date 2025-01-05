package log

import (
  "io"
  "os"
  "testing"

  "github.com/stretchr/testify/require"
)

func TestIndex(t *testing.T) {
  f, err := os.CreateTemp(os.TempDir(), "index_test")
  require.NoError(t,err)
  defer os.Remove(f.Name())

  c :=Config{}
  c.Segment.MaxIndexBytes = 1024
  idx, err := newIndex(f,c)
  require.NoError(t, err)
  _, _, err = idx.Read(-1)
  require.Error(t, err)
  require.Equal(t, f.Name(), idx.Name())
  entries := []struct {
    Off uint32
    Pos uint64
  }{
    {Off: 0, Pos: 0},
    {Off: 1, Pos: 10},
  }
  for _, want := range entries {
    err = idx.Write(want.Off, want.Pos)
    require.NoError(t, err)

    _, pos, err := idx.Read(int64(want.Off))
    require.NoError(t, err)
    require.Equal(t, want.Pos, pos)
  }

  // 존재하는 항목의 범위를 넘어서서 읽으려 하면 에러가 나야 한다.
  _, _, err = idx.Read(int64(len(entries)))
  require.Equal(t, io.EOF, err)
  _ = idx.Close()

  // 파일이 있다면, 파일의 데이터에서 인덱스의 초기 상태를 만들어야 한다.
  f, _ = os.OpenFile(f.Name(), os.O_RDWR, 0600)
  idx, err = newIndex(f, c)
  require.NoError(t, err)
  off, pos, err := idx.Read(-1)
  require.Equal(t, uint32(1), off)
  require.Equal(t, entries[1].Pos, pos)
}
