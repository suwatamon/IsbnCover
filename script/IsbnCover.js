$(function () {

    const PIXEL_SIZE = 28;
    const CANVAS_FACTOR = 5;
    const canvasWidth = PIXEL_SIZE * CANVAS_FACTOR;
    const canvasHeight = canvasWidth;

    function createCanvas() {
        const mem_canvas = document.createElement("canvas");
        mem_canvas.width = canvasWidth;
        mem_canvas.height = canvasHeight;
        mem_canvas.setAttribute("style", "border: solid");
        const context = mem_canvas.getContext('2d');
        context.fillStyle = "rgb(f,f,f)";
        return mem_canvas;
    }

    function createButton() {
        const mem_button = document.createElement("input");
        mem_button.type = "button";
        mem_button.value = "クリア";
        mem_button.addEventListener('click', e => clearCanvas(e.target.canvas));

        return mem_button;
    }

    function createDiv() {
        const mem_div = document.createElement("div");
        mem_div.classList.add("input-container");

        return mem_div;
    }

    function clearCanvas(canvas) {
        const context = canvas.getContext('2d');
        context.fillStyle = "white";
        context.fillRect(0, 0, canvasWidth, canvasHeight);
        context.fillStyle = "black";
        let i = canvas.index;
        for (let j = 0; j < PIXEL_SIZE ** 2; j++)
            pixels[i][j] = 0;
    }

    function setCanvasEvents (canvas, pixel_i){
        canvas.addEventListener("mousedown", e => {
            e.target.drawing = true;
        })

        canvas.addEventListener("mouseup", e => {
            e.target.drawing = false;
        })

        canvas.addEventListener("mousemove", e => {
            const CF = CANVAS_FACTOR;
            const PS = PIXEL_SIZE;
            if (e.target.drawing) {
                let x = Math.floor(e.offsetX / CF);
                let y = Math.floor(e.offsetY / CF);
                if (0 <= x && x < PS && 0 <= y && y < PS) {
                    e.target.getContext("2d").fillRect(x * CF, y * CF, CF, CF);
                    pixel_i[x + y * PS] = 1;
                }
            }
        })
    }

    var canvasList = [];
    var canvas_area = document.getElementById('canvas_area');

    const N_CANVAS = 13;
    for (let i = 0; i < N_CANVAS; i++) {
        const tmpCanvas = createCanvas();
        const tmpButton = createButton();
        const tmpDiv = createDiv();

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
        const canvas = canvasList[i];
        const pixel_i = pixels[i];
        setCanvasEvents(canvas, pixel_i);
    }

    $("#predict").click(function () {
        document.querySelector("#hdnPredict").value = JSON.stringify(pixels);
    });

    $("#all_clear").click(function () {
        for (let i = 0; i < N_CANVAS; i++) {
            clearCanvas(canvasList[i]);
        }
    });

});
