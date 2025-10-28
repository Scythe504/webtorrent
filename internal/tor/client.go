package tor

import (
	"fmt"
	"log"

	"github.com/anacrolix/torrent"
)

type Torrent struct {
	cl  *torrent.Client
	tor map[string]*torrent.Torrent
}

func New() Torrent {
	cfg := torrent.NewDefaultClientConfig()

	cfg.DataDir = "./download/"

	client, err := torrent.NewClient(cfg)

	if err != nil {
		log.Fatal(err)
	}

	return Torrent{
		cl:  client,
		tor: make(map[string]*torrent.Torrent),
	}
}

func (tr *Torrent) AddMagnet(id, magnetLink string) error {
	t, err := tr.cl.AddMagnet(magnetLink)

	if err != nil {
		return err
	}
	<-t.GotInfo()

	if tr.tor[id] != nil {
		return nil
	}
	tr.tor[id] = t

	return nil
}

func (tr *Torrent) GetReader(id string) torrent.Reader {
	t := tr.tor[id]

	if t == nil {
		return nil
	}

	file := tr.tor[id].Files()[0]

	return file.NewReader()
}

func (tr *Torrent) GetMagnetLink(videoId string) *string {
	metainfo := tr.tor[videoId].Metainfo()

	magnetLinkV2, err := metainfo.MagnetV2()

	if err != nil {
		log.Println("[GetMagnetLink]", magnetLinkV2, err)
		return nil
	}

	magnetUri := magnetLinkV2.String()

	return &magnetUri
}

func (tr *Torrent) GetFileName(videoId string) (string, error) {
	t := tr.tor[videoId]
	
	if t == nil {
		return "", fmt.Errorf("torrent not found for videoId: %s", videoId)
	}
	
	if len(t.Files()) == 0 {
		return "", fmt.Errorf("no files in torrent for videoId: %s", videoId)
	}
	
	return t.Files()[0].DisplayPath(), nil
}