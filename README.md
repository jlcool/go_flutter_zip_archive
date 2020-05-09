# GO_FLUTTER_ZIP_ARCHIVE

go-flutter zip 压缩插件，依照我的另一个flutter版插件写的 [flutter_zip_archive](https://github.com/jlcool/flutter_zip_archive)，接口可以在这里看

只能压缩和解压一层文件

## Usage

#### 1. 安装 https://github.com/jlcool/flutter_zip_archive

#### 2. 在cmd/options.go中添加

Import as:

```go
import "github.com/jlcool/go_flutter_zip_archive"
```

```go
flutter.AddPlugin(&ziparchive.ZipArchivePlugin{}),
```
