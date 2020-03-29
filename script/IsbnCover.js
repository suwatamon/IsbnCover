$(function () {

    const PIXEL_SIZE = 28;
    const CANVAS_FACTOR = 5;
    const canvasWidth = PIXEL_SIZE * CANVAS_FACTOR;
    const canvasHeight = canvasWidth;

    function createCanvas() {
        var mem_canvas;
        mem_canvas = document.createElement("canvas");
        mem_canvas.width = canvasWidth;
        mem_canvas.height = canvasHeight;
        mem_canvas.setAttribute("style", "border: solid");
        var context = mem_canvas.getContext('2d');
        context.fillStyle = "rgb(f,f,f)";
        return mem_canvas;
    }

    function createButton() {
        let mem_button;
        mem_button = document.createElement("input");
        mem_button.type = "button";
        mem_button.value = "クリア";
        mem_button.addEventListener('click', e => clearCanvas(e.target.canvas));

        return mem_button;
    }

    function createDiv() {
        let mem_div;
        mem_div = document.createElement("div");
        mem_div.classList.add("input-container");

        return mem_div;
    }


    var canvasList = [];
    var buttonList = [];
    var canvas_area = document.getElementById('canvas_area');

    const N_CANVAS = 13;
    for (var i = 0; i < N_CANVAS; i++) {
        let tmpCanvas = createCanvas();
        let tmpButton = createButton();
        let tmpDiv = createDiv();

        tmpCanvas.index = i;
        canvasList.push(tmpCanvas);
        tmpButton.canvas = tmpCanvas;

        tmpDiv.appendChild(tmpButton);
        tmpDiv.appendChild(tmpCanvas);
        canvas_area.appendChild(tmpDiv);
    }



    var pixels = [];

    for (let j = 0; j < N_CANVAS; j++) {
        pixels[j] = [];
        for (let i = 0; i < PIXEL_SIZE ** 2; i++)
            pixels[j][i] = 0;
    }

    for (let i = 0; i < N_CANVAS; i++) {
        var canvas = canvasList[i];
        // var drawing = false;

        canvas.addEventListener("mousedown", function (e) {
            e.target.drawing = true;
        })

        canvas.addEventListener("mouseup", function (e) {
            e.target.drawing = false;
        })

        canvas.addEventListener("mousemove", function (e) {
            const CF = CANVAS_FACTOR;
            const PS = PIXEL_SIZE;
            if (e.target.drawing) {
                var x = Math.floor(e.offsetX / CF);
                var y = Math.floor(e.offsetY / CF);
                if (0 <= x && x < PS && 0 <= y && y < PS) {
                    e.target.getContext("2d").fillRect(x * CF, y * CF, CF, CF);
                    pixels[i][x + y * PS] = 1;
                }
            }
        })
    }

    $("#predict").click(function () {
        $.ajax({
            url: "http://localhost:8888/predict",
            type: "POST",
            data: {
                "imageList": JSON.stringify(pixels)
            },
            success: function (result) {
                document.write(result);
            }
        })
    });

    $("#all_clear").click(function () {
        for (var i = 0; i < N_CANVAS; i++) {
            clearCanvas(canvasList[i]);
        }
    });

    function clearCanvas(canvas) {
        var context = canvas.getContext('2d');
        context.fillStyle = "white";
        context.fillRect(0, 0, canvasWidth, canvasHeight);
        context.fillStyle = "black";
        let i = canvas.index;
        for (let j = 0; j < PIXEL_SIZE ** 2; j++)
            pixels[i][j] = 0;
    }

});
