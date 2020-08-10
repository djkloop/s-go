package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["fuck_youku_api/controllers:BarrageController"] = append(beego.GlobalControllerRouter["fuck_youku_api/controllers:BarrageController"],
        beego.ControllerComments{
            Method: "Save",
            Router: "/barrage/save",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["fuck_youku_api/controllers:BarrageController"] = append(beego.GlobalControllerRouter["fuck_youku_api/controllers:BarrageController"],
        beego.ControllerComments{
            Method: "BarrageWs",
            Router: "/barrage/ws",
            AllowHTTPMethods: []string{"*"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["fuck_youku_api/controllers:BaseControllers"] = append(beego.GlobalControllerRouter["fuck_youku_api/controllers:BaseControllers"],
        beego.ControllerComments{
            Method: "ChannelRegion",
            Router: "/channel/region",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["fuck_youku_api/controllers:BaseControllers"] = append(beego.GlobalControllerRouter["fuck_youku_api/controllers:BaseControllers"],
        beego.ControllerComments{
            Method: "ChannelType",
            Router: "/channel/type",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["fuck_youku_api/controllers:CommentController"] = append(beego.GlobalControllerRouter["fuck_youku_api/controllers:CommentController"],
        beego.ControllerComments{
            Method: "List",
            Router: "/comment/list",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["fuck_youku_api/controllers:CommentController"] = append(beego.GlobalControllerRouter["fuck_youku_api/controllers:CommentController"],
        beego.ControllerComments{
            Method: "Save",
            Router: "/comment/save",
            AllowHTTPMethods: []string{"*"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["fuck_youku_api/controllers:TopController"] = append(beego.GlobalControllerRouter["fuck_youku_api/controllers:TopController"],
        beego.ControllerComments{
            Method: "ChannelTop",
            Router: "/channel/top",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["fuck_youku_api/controllers:TopController"] = append(beego.GlobalControllerRouter["fuck_youku_api/controllers:TopController"],
        beego.ControllerComments{
            Method: "TypeTop",
            Router: "/type/top",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["fuck_youku_api/controllers:UserController"] = append(beego.GlobalControllerRouter["fuck_youku_api/controllers:UserController"],
        beego.ControllerComments{
            Method: "UserLogin",
            Router: "/login/do",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["fuck_youku_api/controllers:UserController"] = append(beego.GlobalControllerRouter["fuck_youku_api/controllers:UserController"],
        beego.ControllerComments{
            Method: "RegisterUser",
            Router: "/register/save",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["fuck_youku_api/controllers:UserController"] = append(beego.GlobalControllerRouter["fuck_youku_api/controllers:UserController"],
        beego.ControllerComments{
            Method: "SendMessageDo",
            Router: "/send/message",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["fuck_youku_api/controllers:VideoController"] = append(beego.GlobalControllerRouter["fuck_youku_api/controllers:VideoController"],
        beego.ControllerComments{
            Method: "ChannelAdvert",
            Router: "/channel/advert",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["fuck_youku_api/controllers:VideoController"] = append(beego.GlobalControllerRouter["fuck_youku_api/controllers:VideoController"],
        beego.ControllerComments{
            Method: "ChannelHotList",
            Router: "/channel/hot",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["fuck_youku_api/controllers:VideoController"] = append(beego.GlobalControllerRouter["fuck_youku_api/controllers:VideoController"],
        beego.ControllerComments{
            Method: "ChannelRegionRecommendList",
            Router: "/channel/recommend/region",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["fuck_youku_api/controllers:VideoController"] = append(beego.GlobalControllerRouter["fuck_youku_api/controllers:VideoController"],
        beego.ControllerComments{
            Method: "ChannelTypeRecommendList",
            Router: "/channel/recommend/type",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["fuck_youku_api/controllers:VideoController"] = append(beego.GlobalControllerRouter["fuck_youku_api/controllers:VideoController"],
        beego.ControllerComments{
            Method: "ChannelVideo",
            Router: "/channel/video",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["fuck_youku_api/controllers:VideoController"] = append(beego.GlobalControllerRouter["fuck_youku_api/controllers:VideoController"],
        beego.ControllerComments{
            Method: "UserVideo",
            Router: "/user/video",
            AllowHTTPMethods: []string{"*"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["fuck_youku_api/controllers:VideoController"] = append(beego.GlobalControllerRouter["fuck_youku_api/controllers:VideoController"],
        beego.ControllerComments{
            Method: "VideoEpisodesList",
            Router: "/video/episodes/list",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["fuck_youku_api/controllers:VideoController"] = append(beego.GlobalControllerRouter["fuck_youku_api/controllers:VideoController"],
        beego.ControllerComments{
            Method: "VideoInfo",
            Router: "/video/info",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["fuck_youku_api/controllers:VideoController"] = append(beego.GlobalControllerRouter["fuck_youku_api/controllers:VideoController"],
        beego.ControllerComments{
            Method: "VideoSave",
            Router: "/video/save",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["fuck_youku_api/controllers:VideoController"] = append(beego.GlobalControllerRouter["fuck_youku_api/controllers:VideoController"],
        beego.ControllerComments{
            Method: "Search",
            Router: "/video/search",
            AllowHTTPMethods: []string{"*"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["fuck_youku_api/controllers:VideoController"] = append(beego.GlobalControllerRouter["fuck_youku_api/controllers:VideoController"],
        beego.ControllerComments{
            Method: "SendEs",
            Router: "/video/send/es",
            AllowHTTPMethods: []string{"*"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
