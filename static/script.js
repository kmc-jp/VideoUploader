var LastVideoGetTime
var LastCommentGetTime

function editPage(title, user, video, thumbURL) {
    document.getElementById('editHead').textContent = String(title) + "の編集"
    document.getElementById('preview').src = String(thumbURL);
    document.getElementById('editPost').action = "index.up?Page=MyPage&Action=Update&Video=" + String(video)
    document.getElementById('edit').style.display = "unset";
    document.getElementById('video_id').value = String(video);
    ObjectPos('edit');
    return
}


function DeleteConfirm(title, video) {
    document.getElementById('deleteHead').textContent = String(title) + "の削除"
    document.getElementById('delete_button').href = "index.up?Page=MyPage&Action=Delete&Video=" + String(video)
    document.getElementById('delete').style.display = "unset";
    ObjectPos('delete');
    return
}


function SlackMenu(title, user, video, thumbURL, page) {
    document.getElementById('slack_thumb_preview').src = String(thumbURL);
    document.getElementById('slack').style.display = "unset";
    document.getElementById('slack_video_id').value = String(video);
    document.getElementById('slack_video_user').value = String(user);
    document.getElementById('slack_video_page').value = String(page);

    ObjectPos('slack');
    return
}

function ObjectPos(id) {
    let BannerHeightS = document.getElementById(String(id)).clientHeight
    let BannerWidthS = document.getElementById(String(id)).clientWidth

    document.getElementById(String(id)).style.left = String(window.scrollX + (document.documentElement.clientWidth - BannerWidthS) / 2) + "px";
    document.getElementById(String(id)).style.top = String(window.scrollY + (document.documentElement.clientHeight - BannerHeightS) / 2) + "px";

    return
}

function GetScrool() {
    let tgt;
    if ('scrollingElement' in document) {
        tgt = document.scrollingElement;
    } else if (this.browser.isWebKit) {
        tgt = document.body;
    } else {
        tgt = document.documentElement;
    }
    return tgt;
}

function HideWindow(id) {
    document.getElementById(id).style.display = "none";
    return
}


function ToggleDisplayDetails(id) {
    if (document.getElementById(String(id)).style.display == "unset") {
        document.getElementById(String(id)).style.display = "none";
    } else {
        document.getElementById(String(id)).style.display = "unset";
    }
    return
}


function GetSearchInfo() {
    // Action
    var Action = String("")
    for (var i = 0; i < document.search_form.Action.length; i++) {
        if (document.search_form.Action[i].checked) {
            Action = document.search_form.Action[i].value;
        }
    }

    // Keywords
    var Keywords = String()
    Keywords = document.getElementById("video_form").value;

    // Mode
    var SearchMode = String()
    for (var i = 0; i < document.search_form.Mode.length; i++) {
        if (document.search_form.Mode[i].checked) {
            SearchMode = document.search_form.Mode[i].value;
        }
    }

    GetSearchResult(Action, Keywords, SearchMode);
    return;

}

function GetSearchResult(Action, Keywords, SearchMode) {

    // サーバに送信
    var url = location.protocol + "//" + location.hostname + location.pathname
        + "?Action=" + Action
        + "&Keywords=" + Keywords
        + "&Mode=" + SearchMode;

    console.log(url);

    // 検索結果の描画先を作る
    if (document.getElementById("search_result") == null) {
        var Card = document.createElement("div");
        var Title = document.createElement("h5");
        var CardBody = document.createElement("div");

        Card.className = "card";

        Title.className = "card-header";
        Title.textContent = "検索結果";

        CardBody.className = "card-body";
        CardBody.id = "search_result";

        Card.appendChild(Title);
        Card.appendChild(CardBody);
        document.getElementById("search_top").appendChild(Card);
    }

    var request = new XMLHttpRequest();
    request.open('GET', url);
    request.responseType = 'json';
    request.onreadystatechange = function () {
        if (request.readyState != 4) {
            document.getElementById('search_result').textContent = String("検索中です…");
        } else if (request.status != 200) {
            document.getElementById('search_result').textContent = String("検索に失敗しました。通信状態を確認してください。");
        } else {
            var result = request.response;
            if (result == null) {
                document.getElementById('search_result').textContent = String("検索に失敗しました。サーバ上で問題が発生している可能性があります。")
                return
            }
            if (result["status"] != 200) {
                document.getElementById('search_result').textContent = String("検索に失敗しました。サーバ上で問題が発生している可能性があります。")
            } else {
                document.getElementById('search_result').textContent = String("")

                if (result["type"] == "videos") {
                    DisplaySearchVideoResults(result["videos"]);
                    return;
                }
                if (result["type"] == "tags") {
                    DisplaySearchTagResults(result["tags"]);
                    return;
                }
            }
        }
        return;
    }

    request.send(null);
    return;
}

var pre_comments

//コメント内容を取得
function GetComments(VideoID, mode) {
    if (document.getElementById("comment_error") != null) {
        document.getElementById("comment_error").remove();
    }
    var url = location.protocol + "//" + location.hostname + location.pathname + "?Action=GetComment&Video=" + VideoID;

    var request = new XMLHttpRequest();
    request.open('GET', url);
    request.responseType = 'json';
    request.onreadystatechange = function () {
        if ((request.readyState != 4) && (mode == String("init"))) {
            document.getElementById('comment_text').textContent = String("コメントの読み込み中です…")
        } else if ((request.status != 200) && (mode == String("init"))) {
            document.getElementById('comment_text').textContent = String("コメントの取得に失敗しました。")
        } else {
            var result = request.response;
            if (result == null) {
                if (mode == String("init")) {
                    document.getElementById('comment_text').textContent = String("コメントの取得に失敗しました。")
                    return
                } else {
                    return
                }
            }
            if ((result["status"] != 200) && (mode == String("init"))) {
                document.getElementById('comment_text').textContent = String("コメントの取得に失敗しました。")
            } else {
                document.getElementById('comment_text').textContent = String("")
                //Commentの描画
                DisplayComment(result["comments"])
            }
        }
    };
    request.send(null);
}


//コメントを送信する。
function SendComment(VideoID) {
    if (document.getElementById("comment_error") != null) {
        document.getElementById("comment_error").remove();
    }

    var url = location.protocol + "//" + location.hostname + location.pathname + "?Action=AddComment";


    let input = document.getElementById("comment_input").value;

    if (input == "") {
        return
    }

    var data = "Comment=" + input + "&VideoID=" + VideoID;

    document.getElementById("comment_input").value = "";


    var request = new XMLHttpRequest();
    request.open('POST', url);
    request.responseType = "json";
    request.onreadystatechange = function () {
        if (request.readyState != 4) {

        } else if (request.status != 200) {
            var errorMessage = document.createElement.appendChild("div")
            document.getElementById("comment_output").appendChild(errorMessage)
            errorMessage.id = "comment_error"
            errorMessage.textContent = String("コメントの送信に失敗しました。\n通信環境を確認してください。")

        } else {
            var result = request.response;

            if (result["status"] != 200) {
                var errorMessage = document.getElementById("comment_output").appendChild("div")
                errorMessage.id = "comment_error"
                errorMessage.textContent = String("コメントの送信に失敗しました。")
            } else {
                GetComments(VideoID)
            }
        }
    };
    request.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
    request.send(data);
}

function GetVideos(mode) {
    var url = location.protocol + "//" + location.hostname + location.pathname + "?Action=VideoStatus";

    var request = new XMLHttpRequest();
    request.open('GET', url);
    request.responseType = 'json';
    request.onreadystatechange = function () {
        if ((request.readyState != 4) && (mode == String("init"))) {
            document.getElementById('my_videos').textContent = String("Now Loading...")
        } else if ((request.status != 200) && (mode == String("init"))) {
            document.getElementById('my_videos').textContent = String("動画の取得に失敗しました。")
        } else {
            var result = request.response;
            if (result == null) {
                if (mode == String("init")) {
                    document.getElementById('my_videos').textContent = String("動画の取得に失敗しました。")
                    return;
                } else {
                    return;
                }
            }
            if ((result["status"] != 200) && (mode == String("init"))) {
                document.getElementById('my_videos').textContent = String("動画の取得に失敗しました。")
            } else {
                document.getElementById('my_videos').textContent = String("")
                let time = new Date(result["time"])
                if (LastVideoGetTime != null && LastVideoGetTime == time) {
                    return;
                }
                //動画リストの描画
                DisplayVideoList(result["video"]);
                LastVideoGetTime = time;
            }
        }
    };
    request.send(null);
}

//コメントを描画する。
function DisplayComment(comments) {
    if (pre_comments == comments) {
        return
    }
    var CommentDiv = document.getElementById('comment_text')

    if (comments.length == 0) {
        CommentDiv.textContent = String("この動画にはまだコメントがありません。")
        return
    }

    var Table = document.createElement('table');
    Table.className = 'nb';

    CommentDiv.appendChild(Table)

    for (let i = 0; i < comments.length; i++) {
        var Row1 = document.createElement('tr');
        var Row2 = document.createElement('tr');
        var Row3 = document.createElement('tr');

        var Icon = document.createElement('td');
        var ID = document.createElement('td');
        var Text = document.createElement('td');

        var IconImg = document.createElement('img');

        Icon.appendChild(IconImg)

        Row1.appendChild(Icon)
        Row1.appendChild(ID)
        Row2.appendChild(Text)

        Table.appendChild(Row1)
        Table.appendChild(Row2)
        Table.appendChild(Row3)

        // アイコンの設定
        Icon.rowSpan = 3;

        Icon.width = "60px"

        Icon.style.padding = "0px";
        Icon.style.width = "60px"

        IconImg.src = comments[i]['user_icon'];
        IconImg.className = "icon";


        // ユーザ名の埋め込み
        ID.rowSpan = 1;
        ID.textContent = comments[i]['user_name'];
        ID.style.padding = "0px";
        ID.style.textAlign = "left";
        ID.style.width = "unset";

        // 本文の設定
        Text.rowSpan = 2;
        Text.textContent = comments[i]["comment"];
        Text.style.padding = "0px";
        Text.style.textAlign = "left";
        Text.style.width = "unset";

    }

    pre_comments = comments

}

function DisplayVideoList(videos) {
    let VideoList = document.getElementById("my_videos")

    VideoList.textContent = "";
    if (videos == null || videos.length == 0) {
        VideoList.textContent = "まだあなたは動画を投稿していないようです。"
    }

    let Table = document.createElement("table");
    Table.className = "nb";

    // エラー
    for (let i = 0; i < videos.length; i++) {
        if (videos[i]["error"] != "") {
            let TopRow = document.createElement("tr");
            // 一行目
            {
                time = new Date(videos[i]["time"])
                let TimeText = document.createElement("small");
                TimeText.className = "text-muted";
                TimeText.textContent = "投稿日時："
                    + String(time.getFullYear()) + "年"
                    + String(time.getMonth() + 1) + "月"
                    + String(time.getDate()) + "日"
                    + String(time.getHours()) + "時"
                    + String(time.getMinutes()) + "分";

                let ImgData = document.createElement("td");
                let TimeData = document.createElement("td");

                let Img = document.createElement("img");

                ImgData.rowSpan = 4;
                TimeData.rowSpan = 1;

                Img.src = "static/icon/movie.png";
                Img.className = "video";

                ImgData.appendChild(Img);
                TimeData.appendChild(TimeText);

                TopRow.appendChild(ImgData);
                TopRow.appendChild(TimeData);
            }

            // 二行目
            let MidRow = document.createElement("tr");

            {
                let TitleData = document.createElement("td");
                let Title = document.createElement("h6");
                let TitleLink = document.createElement("a");

                Title.textContent = videos[i]["title"];
                TitleData.rowSpan = 2;

                TitleData.appendChild(Title);
                MidRow.appendChild(TitleData);
            }

            // 三行目
            let BottomRow = document.createElement("tr")
            {
                let EditData = document.createElement("td");
                let SlackData = document.createElement("td");
                let DeleteData = document.createElement("td");

                let EditIcon = document.createElement("img");
                let SlackSpan = document.createElement("span");
                let SlackIcon = document.createElement("i");
                let DeleteIcon = document.createElement("img");

                EditIcon.src = "static/icon/edit.png";
                EditIcon.style.height = "30px";

                SlackSpan.style.fontSize = "30px";
                SlackSpan.style.color = "#4B4B4B";

                SlackIcon.className = "fab fa-slack";

                DeleteIcon.src = "static/icon/delete.png";
                DeleteIcon.style.height = "30px"

                EditData.rowSpan = 2;
                SlackData.rowSpan = 2;
                DeleteData.rowSpan = 2;

                SlackSpan.appendChild(SlackIcon);

                EditData.appendChild(EditIcon);
                SlackData.appendChild(SlackSpan);;
                DeleteData.appendChild(DeleteIcon);

                BottomRow.appendChild(EditData);
                BottomRow.appendChild(SlackData);
                BottomRow.appendChild(DeleteData);
            }

            let StatusRow = document.createElement("tr");
            {
                let IconData = document.createElement("td");
                let TextData = document.createElement("td");
                let StateData = document.createElement("td");

                let IconImg = document.createElement("img");

                let TextP = document.createElement("p");
                let StateP = document.createElement("p");

                TextP.textContent = "エラーが発生しました。";
                StateP.textContent = "現在:" + videos[i]["error"];

                IconImg.src = "static/icon/caution.png";
                IconImg.className = "video";

                IconImg.style.height = "30px";
                IconImg.style.width = "30px";

                IconData.appendChild(IconImg);
                TextData.appendChild(TextP);
                StateData.appendChild(StateP);

                StatusRow.appendChild(IconData);
                StatusRow.appendChild(TextData);
                StatusRow.appendChild(StateData);
                StatusRow.appendChild(TextData2);
            }
            StatusRow.className = "status"

            Table.appendChild(TopRow);
            Table.appendChild(MidRow);
            Table.appendChild(BottomRow);
            Table.appendChild(document.createElement("tr"));
            Table.appendChild(StatusRow);
        }
        if (videos[i]["phase"] != "") {
            let TopRow = document.createElement("tr");
            // 一行目
            {
                time = new Date(videos[i]["time"])
                let TimeText = document.createElement("small");
                TimeText.className = "text-muted";
                TimeText.textContent = "投稿日時："
                    + String(time.getFullYear()) + "年"
                    + String(time.getMonth() + 1) + "月"
                    + String(time.getDate()) + "日"
                    + String(time.getHours()) + "時"
                    + String(time.getMinutes()) + "分";

                let ImgData = document.createElement("td");
                let TimeData = document.createElement("td");

                let Img = document.createElement("img");

                ImgData.rowSpan = 4;
                TimeData.rowSpan = 1;

                Img.src = Img.src = "static/icon/movie.png";;
                Img.className = "video";

                ImgData.appendChild(Img);
                TimeData.appendChild(TimeText);

                TopRow.appendChild(ImgData);
                TopRow.appendChild(TimeData);
            }

            // 二行目
            let MidRow = document.createElement("tr");

            {
                let TitleData = document.createElement("td");
                let Title = document.createElement("h6");

                Title.textContent = videos[i]["title"];
                TitleData.rowSpan = 2;

                TitleData.appendChild(Title);
                MidRow.appendChild(TitleData);
            }

            // 三行目
            let BottomRow = document.createElement("tr")
            {
                let EditData = document.createElement("td");
                let SlackData = document.createElement("td");
                let DeleteData = document.createElement("td");

                let EditIcon = document.createElement("img");
                let SlackSpan = document.createElement("span");
                let SlackIcon = document.createElement("i");
                let DeleteIcon = document.createElement("img");

                EditIcon.src = "static/icon/edit.png";
                EditIcon.style.height = "30px";

                SlackSpan.style.fontSize = "30px";
                SlackSpan.style.color = "#4B4B4B";

                SlackIcon.className = "fab fa-slack";

                DeleteIcon.src = "static/icon/delete.png";
                DeleteIcon.style.height = "30px"

                EditData.rowSpan = 2;
                SlackData.rowSpan = 2;
                DeleteData.rowSpan = 2;

                SlackSpan.appendChild(SlackIcon);

                EditData.appendChild(EditIcon);
                SlackData.appendChild(SlackSpan);;
                DeleteData.appendChild(DeleteIcon);

                BottomRow.appendChild(EditData);
                BottomRow.appendChild(SlackData);
                BottomRow.appendChild(DeleteData);
            }

            let StatusRow = document.createElement("tr");
            {
                let IconData = document.createElement("td");
                let TextData = document.createElement("td");
                let StateData = document.createElement("td");
                let TextData2 = document.createElement("td");

                let IconImg = document.createElement("img");

                let TextP = document.createElement("p");
                let StateP = document.createElement("p");
                let Text2P = document.createElement("p");

                TextP.textContent = "処理中…";
                StateP.textContent = "現在:" + videos[i]["phase"];
                Text2P.textContent = "ステータスが変化しない場合は再読み込みしてください。";

                IconImg.src = "static/icon/loading.gif";
                IconImg.className = "video";

                IconImg.style.height = "30px";
                IconImg.style.width = "30px";

                IconData.appendChild(IconImg);
                TextData.appendChild(TextP);
                StateData.appendChild(StateP);
                TextData2.appendChild(Text2P);

                StatusRow.appendChild(IconData);
                StatusRow.appendChild(TextData);
                StatusRow.appendChild(StateData);
                StatusRow.appendChild(TextData2);
            }
            StatusRow.className = "status"

            Table.appendChild(TopRow);
            Table.appendChild(MidRow);
            Table.appendChild(BottomRow);
            Table.appendChild(document.createElement("tr"));
            Table.appendChild(StatusRow);
        }
        // 通常
        if ((videos[i]["phase"] == "") && (videos[i]["error"] == "")) {
            let TopRow = document.createElement("tr");
            // 一行目
            {
                time = new Date(videos[i]["time"])
                let TimeText = document.createElement("small");
                TimeText.className = "text-muted";
                TimeText.textContent = "投稿日時："
                    + String(time.getFullYear()) + "年"
                    + String(time.getMonth() + 1) + "月"
                    + String(time.getDate()) + "日"
                    + String(time.getHours()) + "時"
                    + String(time.getMinutes()) + "分";

                let ImgData = document.createElement("td");
                let TimeData = document.createElement("td");

                let TimeVideoLink = document.createElement("a");
                let ImgVideoLink = document.createElement("a");

                let Img = document.createElement("img");

                ImgData.rowSpan = 4;
                TimeData.rowSpan = 1;

                ImgVideoLink.href = "index.up?Page=Play&User=" + videos[i]["user"] + "&Video=" + videos[i]["id"];
                TimeVideoLink.href = "index.up?Page=Play&User=" + videos[i]["user"] + "&Video=" + videos[i]["id"];

                Img.src = videos[i]["thumb_url"];
                Img.className = "video";

                ImgVideoLink.appendChild(Img);
                ImgData.appendChild(ImgVideoLink);

                TimeVideoLink.appendChild(TimeText);
                TimeData.appendChild(TimeVideoLink);

                TopRow.appendChild(ImgData);
                TopRow.appendChild(TimeData);
            }

            // 二行目
            let MidRow = document.createElement("tr");

            {
                let TitleData = document.createElement("td");
                let Title = document.createElement("h6");
                let TitleLink = document.createElement("a");

                Title.textContent = videos[i]["title"];

                TitleLink.appendChild(Title);
                TitleLink.href = "index.up?Page=Play&User=" + videos[i]["user"] + "&Video=" + videos[i]["id"];

                TitleData.rowSpan = 2;
                TitleData.appendChild(TitleLink);

                MidRow.appendChild(TitleData);
            }

            // 三行目
            let BottomRow = document.createElement("tr")
            {
                let EditData = document.createElement("td");
                let SlackData = document.createElement("td");
                let DeleteData = document.createElement("td");

                let EditLink = document.createElement("a");
                let DeleteLink = document.createElement("a");
                let SlackLink = document.createElement("a");

                let EditIcon = document.createElement("img");
                let SlackSpan = document.createElement("span");
                let SlackIcon = document.createElement("i");
                let DeleteIcon = document.createElement("img");

                EditLink.href = "javascript:editPage('"
                    + videos[i]["title"] + "','"
                    + videos[i]["user"] + "','"
                    + videos[i]["id"] + "','"
                    + videos[i]["thumb_url"] + "');";

                EditIcon.src = "static/icon/edit.png";
                EditIcon.style.height = "30px";

                SlackLink.href = "javascript:SlackMenu('"
                    + videos[i]["title"] + "','"
                    + videos[i]["user"] + "','"
                    + videos[i]["id"] + "','"
                    + videos[i]["thumb_url"]
                    + "', 'MyPage');";

                SlackSpan.style.fontSize = "30px";
                SlackSpan.style.color = "#4B4B4B";

                SlackIcon.className = "fab fa-slack";

                DeleteLink.href = "javascript:DeleteConfirm('"
                    + videos[i]["title"] + "','"
                    + videos[i]["id"] + "');";

                DeleteIcon.src = "static/icon/delete.png";
                DeleteIcon.style.height = "30px"

                EditData.rowSpan = 2;
                SlackData.rowSpan = 2;
                DeleteData.rowSpan = 2;


                EditLink.appendChild(EditIcon);

                SlackSpan.appendChild(SlackIcon);
                SlackLink.appendChild(SlackSpan);

                DeleteLink.appendChild(DeleteIcon);

                EditData.appendChild(EditLink);
                SlackData.appendChild(SlackLink);
                DeleteData.appendChild(DeleteLink);

                BottomRow.appendChild(EditData);
                BottomRow.appendChild(SlackData);
                BottomRow.appendChild(DeleteData);
            }

            Table.appendChild(TopRow);
            Table.appendChild(MidRow);
            Table.appendChild(BottomRow);
            Table.appendChild(document.createElement("tr"));
        }
    }
    VideoList.appendChild(Table)

    return
}


// 検索結果描画
function DisplaySearchVideoResults(videos) {
    let VideoList = document.getElementById("search_result")
    VideoList.textContent = "";
    if (videos == null || videos.length == 0) {
        VideoList.textContent = "該当する動画は見つかりませんでした。";
    }
    let Table = document.createElement("table");
    Table.className = "nb";

    for (let i = 0; i < videos.length; i++) {

        if (videos[i]["phase"] == "" && videos[i]["error"] == "") {
            let TopRow = document.createElement("tr");
            // 一行目
            {
                time = new Date(videos[i]["time"])
                let TimeText = document.createElement("small");
                TimeText.className = "text-muted";
                TimeText.textContent = "投稿日時："
                    + String(time.getFullYear()) + "年"
                    + String(time.getMonth() + 1) + "月"
                    + String(time.getDate()) + "日"
                    + String(time.getHours()) + "時"
                    + String(time.getMinutes()) + "分";

                let ImgData = document.createElement("td");
                let TimeData = document.createElement("td");

                let TimeVideoLink = document.createElement("a");
                let ImgVideoLink = document.createElement("a");

                let Img = document.createElement("img");

                ImgData.rowSpan = 4;
                TimeData.rowSpan = 1;

                ImgVideoLink.href = "index.up?Page=Play&User=" + videos[i]["user"] + "&Video=" + videos[i]["id"];
                TimeVideoLink.href = "index.up?Page=Play&User=" + videos[i]["user"] + "&Video=" + videos[i]["id"];

                Img.src = videos[i]["thumb_url"];
                Img.className = "video";

                ImgVideoLink.appendChild(Img);
                ImgData.appendChild(ImgVideoLink);

                TimeVideoLink.appendChild(TimeText);
                TimeData.appendChild(TimeVideoLink);

                TopRow.appendChild(ImgData);
                TopRow.appendChild(TimeData);
            }

            // 二行目
            let MidRow = document.createElement("tr");

            {
                let TitleData = document.createElement("td");
                let Title = document.createElement("h6");
                let TitleLink = document.createElement("a");

                Title.textContent = videos[i]["title"];

                TitleLink.appendChild(Title);
                TitleLink.href = "index.up?Page=Play&User=" + videos[i]["user"] + "&Video=" + videos[i]["id"];

                TitleData.rowSpan = 2;
                TitleData.appendChild(TitleLink);

                MidRow.appendChild(TitleData);
            }

            // 三行目
            let BottomRow = document.createElement("tr")
            {
                let SlackData = document.createElement("td");
                let SlackLink = document.createElement("a");
                let SlackSpan = document.createElement("span");
                let SlackIcon = document.createElement("i");

                SlackLink.href = "javascript:SlackMenu('"
                    + videos[i]["title"] + "','"
                    + videos[i]["user"] + "','"
                    + videos[i]["id"] + "','"
                    + videos[i]["thumb_url"]
                    + "', 'MyPage');";

                SlackSpan.style.fontSize = "30px";
                SlackSpan.style.color = "#4B4B4B";

                SlackIcon.className = "fab fa-slack";
                SlackData.rowSpan = 2;

                SlackSpan.appendChild(SlackIcon);
                SlackLink.appendChild(SlackSpan);
                SlackData.appendChild(SlackLink);

                BottomRow.appendChild(SlackData);
            }

            Table.appendChild(TopRow);
            Table.appendChild(MidRow);
            Table.appendChild(BottomRow);
            Table.appendChild(document.createElement("tr"));

            VideoList.appendChild(Table)
        }

        return
    }
}

function DisplaySearchTagResults(tags) {
    let TagList = document.getElementById("search_result")
    TagList.textContent = "";
    if (tags == null || tags.length == 0) {
        TagList.textContent = "該当する動画は見つかりませんでした。";
    }
    let Table = document.createElement("table");
    Table.className = "nb";

    for (let i = 0; i < tags.length; i++) {
        var Row = document.createElement("tr");
        {
            var IconData = document.createElement("td");
            var IconSpan = document.createElement("span");
            var Icon = document.createElement("i");

            IconData.style.width = "60px";
            IconData.style.padding = "0px";

            IconSpan.style.fontSize = "60px";

            Icon.className = "fas fa-tag";

            IconSpan.appendChild(Icon);
            IconData.appendChild(IconSpan);
            Row.appendChild(IconData);
        }
        {
            var TagData = document.createElement("td");
            var TagLink = document.createElement("a")

            TagData.style.textAlign = "left";
            TagData.style.width = "unset";

            TagLink.href = "javascript:GetSearchResult('SearchTags'," + "'" + tags[i] + "','or')";
            TagLink.textContent = tags[i];

            TagData.appendChild(TagLink);

            Row.appendChild(TagData)
        }
        Table.appendChild(Row);
    }

    TagList.appendChild(Table);
    return;
}
