package server

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/photoprism/photoprism/internal/api"
	"github.com/photoprism/photoprism/internal/config"
)

func registerRoutes(router *gin.Engine, conf *config.Config) {
	// Enables automatic redirection if the current route cannot be matched but a
	// handler for the path with (without) the trailing slash exists.
	router.RedirectTrailingSlash = true

	// Static favicon image.
	router.StaticFile(conf.BaseUri("/favicon.ico"), filepath.Join(conf.ImgPath(), "favicon.ico"))

	// Static assets like js, css and font files.
	router.Static(conf.BaseUri(config.StaticUri), conf.StaticPath())

	// Serve custom static assets if dir exists.
	if dir := conf.CustomStaticPath(); dir != "" {
		router.Static(conf.BaseUri(config.CustomStaticUri), dir)
	}

	// PWA Manifest.
	router.GET(conf.BaseUri("/manifest.json"), func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("Content-Type", "application/json")

		clientConfig := conf.ClientPublic()
		c.HTML(http.StatusOK, "manifest.json", gin.H{"config": clientConfig})
	})

	// PWA Service Worker.
	router.GET(conf.BaseUri("/sw.js"), func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.File(filepath.Join(conf.BuildPath(), "sw.js"))
	})

	// Rainbow Page.
	router.GET(conf.BaseUri("/_rainbow"), func(c *gin.Context) {
		clientConfig := conf.ClientPublic()
		c.HTML(http.StatusOK, "rainbow.tmpl", gin.H{"config": clientConfig})
	})

	// Splash Screen.
	router.GET(conf.BaseUri("/_splash"), func(c *gin.Context) {
		clientConfig := conf.ClientPublic()
		c.HTML(http.StatusOK, "splash.tmpl", gin.H{"config": clientConfig})
	})

	// JSON-REST API Version 1
	v1 := router.Group(conf.BaseUri(config.ApiUri))
	{
		// Global App Config.
		api.GetConfigOptions(v1)
		api.SaveConfigOptions(v1)

		// Custom Settings.
		api.GetClientConfig(v1)
		api.GetSettings(v1)
		api.SaveSettings(v1)

		// User Authentication.
		api.ChangePassword(v1)
		api.CreateSession(v1)
		api.GetSession(v1)
		api.DeleteSession(v1)

		// External Account Management.
		api.SearchServices(v1)
		api.GetService(v1)
		api.GetServiceFolders(v1)
		api.UploadToService(v1)
		api.AddService(v1)
		api.DeleteService(v1)
		api.UpdateService(v1)

		// Thumbs, Downloads, and Videos.
		api.GetThumb(v1)
		api.GetDownload(v1)
		api.GetVideo(v1)
		api.ZipCreate(v1)
		api.ZipDownload(v1)

		// Index, Upload, and Import.
		api.Upload(v1)
		api.StartImport(v1)
		api.CancelImport(v1)
		api.StartIndexing(v1)
		api.CancelIndexing(v1)

		// Search & Organization.
		api.SearchPhotos(v1)
		api.SearchGeo(v1)
		api.GetPhoto(v1)
		api.GetPhotoYaml(v1)
		api.UpdatePhoto(v1)
		api.GetPhotoDownload(v1)
		// api.GetPhotoLinks(v1)
		// api.CreatePhotoLink(v1)
		// api.UpdatePhotoLink(v1)
		// api.DeletePhotoLink(v1)
		api.ApprovePhoto(v1)
		api.LikePhoto(v1)
		api.DislikePhoto(v1)
		api.AddPhotoLabel(v1)
		api.RemovePhotoLabel(v1)
		api.UpdatePhotoLabel(v1)
		api.GetMomentsTime(v1)
		api.GetFile(v1)
		api.DeleteFile(v1)
		api.UpdateMarker(v1)
		api.ClearMarkerSubject(v1)
		api.PhotoPrimary(v1)
		api.PhotoUnstack(v1)

		// Photo Albums.
		api.SearchAlbums(v1)
		api.GetAlbum(v1)
		api.AlbumCover(v1)
		api.CreateAlbum(v1)
		api.UpdateAlbum(v1)
		api.DeleteAlbum(v1)
		api.DownloadAlbum(v1)
		api.GetAlbumLinks(v1)
		api.CreateAlbumLink(v1)
		api.UpdateAlbumLink(v1)
		api.DeleteAlbumLink(v1)
		api.LikeAlbum(v1)
		api.DislikeAlbum(v1)
		api.CloneAlbums(v1)
		api.AddPhotosToAlbum(v1)
		api.RemovePhotosFromAlbum(v1)

		// Photo Labels.
		api.SearchLabels(v1)
		api.LabelCover(v1)
		api.UpdateLabel(v1)
		// api.GetLabelLinks(v1)
		// api.CreateLabelLink(v1)
		// api.UpdateLabelLink(v1)
		// api.DeleteLabelLink(v1)
		api.LikeLabel(v1)
		api.DislikeLabel(v1)

		// Files and Folders.
		api.SearchFoldersOriginals(v1)
		api.SearchFoldersImport(v1)
		api.FolderCover(v1)

		// People.
		api.SearchSubjects(v1)
		api.GetSubject(v1)
		api.UpdateSubject(v1)
		api.LikeSubject(v1)
		api.DislikeSubject(v1)

		// Faces.
		api.SearchFaces(v1)
		api.GetFace(v1)
		api.UpdateFace(v1)

		// Batch Operations.
		api.BatchPhotosApprove(v1)
		api.BatchPhotosArchive(v1)
		api.BatchPhotosRestore(v1)
		api.BatchPhotosPrivate(v1)
		api.BatchPhotosDelete(v1)
		api.BatchAlbumsDelete(v1)
		api.BatchLabelsDelete(v1)

		// Technical Endpoints.
		api.GetSvg(v1)
		api.GetStatus(v1)
		api.GetErrors(v1)
		api.DeleteErrors(v1)
		api.SendFeedback(v1)
		api.Connect(v1)
		api.WebSocket(v1)
	}

	// Configure link sharing.
	s := router.Group(conf.BaseUri("/s"))
	{
		api.Shares(s)
		api.SharePreview(s)
	}

	// WebDAV server for file management, sync and sharing.
	if conf.DisableWebDAV() {
		log.Info("webdav: server disabled")
	} else {
		WebDAV(conf.OriginalsPath(), router.Group(conf.BaseUri(WebDAVOriginals), BasicAuth()), conf)
		log.Infof("webdav: %s/ enabled, waiting for requests", conf.BaseUri(WebDAVOriginals))

		if conf.ImportPath() != "" {
			WebDAV(conf.ImportPath(), router.Group(conf.BaseUri(WebDAVImport), BasicAuth()), conf)
			log.Infof("webdav: %s/ enabled, waiting for requests", conf.BaseUri(WebDAVImport))
		}
	}

	// User Interface.
	router.GET(conf.BaseUri("/library/*path"), func(c *gin.Context) {
		values := gin.H{
			"signUp": gin.H{"message": config.MsgSponsor, "url": config.SignUpURL},
			"config": conf.ClientPublic(),
		}
		c.HTML(http.StatusOK, conf.TemplateName(), values)
	})

	// Redirect to UI.
	router.GET(conf.BaseUri("/"), func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, conf.BaseUri("/library/login"))
	})
}
