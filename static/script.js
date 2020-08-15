
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

function HideEditWindow() {
    document.getElementById('edit').style.display = "none";
    return
}


function HideSlackMenu() {
    document.getElementById('slack').style.display = "none";
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

function HideDeleteWindow() {
    document.getElementById('delete').style.display = "none";
    return
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

var pre_video_list

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
                    return
                } else {
                    return
                }
            }
            if ((result["status"] != 200) && (mode == String("init"))) {
                document.getElementById('my_videos').textContent = String("動画の取得に失敗しました。")
            } else {
                document.getElementById('my_videos').textContent = String("")
                //動画リストの描画
                DisplayVideoList(result["video"])
            }
        }
    };
    request.send(null);
}


function DisplayVideoList(videos) {
    if (pre_video_list == videos) {
        return
    }

    let VideoList = document.getElementById("my_videos")

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

                ImgVideoLink.href = videos[i]["url"];
                TimeVideoLink.href = videos[i]["url"];

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
                TitleLink.href = videos[i]["url"];

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

    pre_video_list = videos

    return
}