package plugin

import (
	. "github.com/iilaurens/filestash/server/common"
	_ "github.com/iilaurens/filestash/server/plugin/plg_authenticate_admin"
	_ "github.com/iilaurens/filestash/server/plugin/plg_authenticate_htpasswd"
	_ "github.com/iilaurens/filestash/server/plugin/plg_authenticate_ldap"
	_ "github.com/iilaurens/filestash/server/plugin/plg_authenticate_openid"
	_ "github.com/iilaurens/filestash/server/plugin/plg_authenticate_passthrough"
	_ "github.com/iilaurens/filestash/server/plugin/plg_authenticate_saml"
	_ "github.com/iilaurens/filestash/server/plugin/plg_backend_artifactory"
	_ "github.com/iilaurens/filestash/server/plugin/plg_backend_backblaze"
	_ "github.com/iilaurens/filestash/server/plugin/plg_backend_dav"
	_ "github.com/iilaurens/filestash/server/plugin/plg_backend_dropbox"
	_ "github.com/iilaurens/filestash/server/plugin/plg_backend_ftp"
	_ "github.com/iilaurens/filestash/server/plugin/plg_backend_gdrive"
	_ "github.com/iilaurens/filestash/server/plugin/plg_backend_git"
	_ "github.com/iilaurens/filestash/server/plugin/plg_backend_ldap"
	_ "github.com/iilaurens/filestash/server/plugin/plg_backend_local"
	_ "github.com/iilaurens/filestash/server/plugin/plg_backend_mysql"
	_ "github.com/iilaurens/filestash/server/plugin/plg_backend_nfs"
	_ "github.com/iilaurens/filestash/server/plugin/plg_backend_nop"
	_ "github.com/iilaurens/filestash/server/plugin/plg_backend_s3"
	_ "github.com/iilaurens/filestash/server/plugin/plg_backend_samba"
	_ "github.com/iilaurens/filestash/server/plugin/plg_backend_sftp"
	_ "github.com/iilaurens/filestash/server/plugin/plg_backend_storj"
	_ "github.com/iilaurens/filestash/server/plugin/plg_backend_tmp"
	_ "github.com/iilaurens/filestash/server/plugin/plg_backend_webdav"
	_ "github.com/iilaurens/filestash/server/plugin/plg_editor_onlyoffice"
	_ "github.com/iilaurens/filestash/server/plugin/plg_handler_console"
	_ "github.com/iilaurens/filestash/server/plugin/plg_image_ascii"
	_ "github.com/iilaurens/filestash/server/plugin/plg_image_thumbnail"
	_ "github.com/iilaurens/filestash/server/plugin/plg_image_transcode"
	_ "github.com/iilaurens/filestash/server/plugin/plg_search_stateless"
	_ "github.com/iilaurens/filestash/server/plugin/plg_security_scanner"
	_ "github.com/iilaurens/filestash/server/plugin/plg_security_svg"
	_ "github.com/iilaurens/filestash/server/plugin/plg_starter_http"
	_ "github.com/iilaurens/filestash/server/plugin/plg_video_thumbnail"
	_ "github.com/iilaurens/filestash/server/plugin/plg_video_transcoder"
)

func init() {
	Log.Debug("Plugin loader")
}
