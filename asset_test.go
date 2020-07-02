package kasset_test

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"git.kanosolution.net/kano/appkit"
	"git.kanosolution.net/kano/kaos"
	"github.com/eaciit/toolkit"
	cv "github.com/smartystreets/goconvey/convey"

	hc "git.kanosolution.net/kano/kaos/client-engines/http"
	hd "git.kanosolution.net/kano/kaos/deployer-engines/http"
	_ "github.com/ariefdarmawan/flexpg"
	"github.com/ariefdarmawan/kasset"
)

func TestAsset(t *testing.T) {
	var (
		e           error
		serviceName = "asset"
		version     = "v1"
		log         = appkit.Log()
		basePath    = os.TempDir()
		hostName    = "localhost:8097"
	)

	log.LogToStdOut = false
	cv.Convey("preparing server", t, func() {
		defer func() {
			os.Remove(basePath)
		}()

		cv.So(e, cv.ShouldBeNil)

		// service
		svc := kaos.NewService().
			SetBasePoint("api/" + version + "/" + serviceName).
			SetLogger(log)

		// register model
		eng := kasset.NewAssetEngine(kasset.NewSimpleFS(basePath))
		svc.RegisterModel(eng, "assetengine")

		// deployer
		mux := http.NewServeMux()
		e = hd.NewHttpDeployer().Deploy(svc, mux)
		cv.So(e, cv.ShouldBeNil)

		go func() {
			svc.Log().Infof("Running service on %s", hostName)
			http.ListenAndServe(hostName, mux)
		}()

		cv.Convey("reading an asset from physical to be saved", func() {
			cwd, _ := os.Getwd()
			sampleFolder := filepath.Join(cwd, "sample")
			filename := "kano-18.jpg"

			bs, err := ioutil.ReadFile(filepath.Join(sampleFolder, filename))
			cv.So(err, cv.ShouldBeNil)
			cv.So(len(bs), cv.ShouldBeGreaterThan, 0)
			cv.Printf("\nReading %.2f MB of data", float64(len(bs))/1024.0/1024.0)
			c, e := hc.NewHttpClient(hostName, nil)
			cv.So(e, cv.ShouldBeNil)

			cv.Convey("save the asset using different name", func() {
				asset := new(kasset.Asset)
				asset.OriginalFileName = filename
				req := new(kasset.AssetData)
				req.Content = bs
				req.Asset = asset

				err := c.CallTo("/api/v1/asset/assetengine/write", asset, req)
				cv.So(err, cv.ShouldBeNil)
				cv.So(asset.ID, cv.ShouldNotEqual, "")
				cv.Printf("\nAsset: %s\n", toolkit.JsonString(asset))

				cv.Convey("read the file", func() {
					readResult := new(kasset.AssetData)
					err := c.CallTo("/api/v1/asset/assetengine/read", readResult, asset.ID)
					cv.So(err, cv.ShouldBeNil)
					cv.So(readResult.Asset.ID, cv.ShouldEqual, asset.ID)
					cv.So(len(readResult.Content), cv.ShouldEqual, len(bs))
					cv.Printf("\nAsset: %s\n", toolkit.JsonString(readResult.Asset))
					cv.Printf("Content: %d\n", len(readResult.Content))

					cv.Convey("delete the file", func() {
						result := int(0)
						err := c.CallTo("/api/v1/asset/assetengine/delete", &result, asset.ID)
						cv.So(err, cv.ShouldBeNil)
						cv.So(result, cv.ShouldEqual, len(readResult.Content))

						cv.Convey("read again the file, should be EOF", func() {
							readResultEOF := new(kasset.AssetData)
							err := c.CallTo("/api/v1/asset/assetengine/read", readResultEOF, asset.ID)
							cv.So(err, cv.ShouldNotBeNil)
						})
					})
				})
			})

			cv.Convey("save asset using target file name", func() {
				asset := new(kasset.Asset)
				asset.OriginalFileName = filename
				asset.NewFileName = filename
				req := new(kasset.AssetData)
				req.Content = bs
				req.Asset = asset

				err := c.CallTo("/api/v1/asset/assetengine/write", asset, req)
				cv.So(err, cv.ShouldBeNil)
				cv.So(asset.ID, cv.ShouldNotEqual, "")
				cv.Printf("\nAsset: %s\n", toolkit.JsonString(asset))

				cv.Convey("read the file", func() {
					readResult := new(kasset.AssetData)
					err := c.CallTo("/api/v1/asset/assetengine/read", readResult, asset.ID)
					cv.So(err, cv.ShouldBeNil)
					cv.So(readResult.Asset.ID, cv.ShouldEqual, asset.ID)
					cv.So(len(readResult.Content), cv.ShouldEqual, len(bs))
					cv.Printf("\nAsset: %s\n", toolkit.JsonString(readResult.Asset))
					cv.Printf("Content: %d\n", len(readResult.Content))
				})
			})
		})
	})
}
