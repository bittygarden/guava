// 棋盘和棋子的配置
const TOP_HEIGHT = 80;  // 顶部区域高度
const BOARD_WIDTH = 900;
const BOARD_HEIGHT = 660 + TOP_HEIGHT;  // 增加总高度
const GRID_SIZE = 60;
const PIECE_RADIUS = 27;
const SIDE_WIDTH = 150;   // 从100改为150

// 获取Canvas上下文
const canvas = document.getElementById('chessboard');
const ctx = canvas.getContext('2d');

// 定义棋子初始位置
const initialBoard = [
    ['車', '馬', '相', '仕', '帥', '仕', '相', '馬', '車'],
    ['', '', '', '', '', '', '', '', ''],
    ['', '炮', '', '', '', '', '', '炮', ''],
    ['兵', '', '兵', '', '兵', '', '兵', '', '兵'],
    ['', '', '', '', '', '', '', '', ''],
    ['', '', '', '', '', '', '', '', ''],
    ['卒', '', '卒', '', '卒', '', '卒', '', '卒'],
    ['', '砲', '', '', '', '', '', '砲', ''],
    ['', '', '', '', '', '', '', '', ''],
    ['俥', '傌', '象', '士', '將', '士', '象', '傌', '俥']
];

let gameBoard = JSON.parse(JSON.stringify(initialBoard));
let selectedPiece = null;
let currentPlayer = 'red'; // 'red' 或 'black'
let hoveredPosition = null;
let possibleMoves = [];

// 添加存储被吃掉棋子的数组
let capturedRed = [];    // 存储被吃掉的红色棋子
let capturedBlack = [];  // 存储被吃掉的黑色棋子

// 在全局变量区域添加
let gameOver = false;
let winner = null;

// 在全局变量区域添加计时器相关变量
let moveStartTime = Date.now();
let redTotalTime = 0;
let blackTotalTime = 0;

// 在全局变量区域添加回合用时变量
let currentTurnStartTime = Date.now();  // 当前回合开始时间
let currentTurnTime = 0;  // 当前回合已用时间

// 绘制棋盘
function drawBoard() {
    ctx.clearRect(0, 0, BOARD_WIDTH, BOARD_HEIGHT);
    
    // 绘制整体背景
    ctx.fillStyle = '#f0d5a8';
    ctx.fillRect(0, 0, BOARD_WIDTH, BOARD_HEIGHT);
    
    // 绘制顶部区域背景
    ctx.fillStyle = '#e0c598';
    ctx.fillRect(0, 0, BOARD_WIDTH, TOP_HEIGHT);
    
    // 绘制两侧区域背景
    ctx.fillStyle = '#e0c598';
    ctx.fillRect(0, TOP_HEIGHT, SIDE_WIDTH, BOARD_HEIGHT - TOP_HEIGHT); // 左侧区域
    ctx.fillRect(BOARD_WIDTH - SIDE_WIDTH, TOP_HEIGHT, SIDE_WIDTH, BOARD_HEIGHT - TOP_HEIGHT); // 右侧区域
    
    // 绘制分隔线
    ctx.strokeStyle = '#000';
    ctx.lineWidth = 2;
    // 顶部区域分隔线
    ctx.beginPath();
    ctx.moveTo(0, TOP_HEIGHT);
    ctx.lineTo(BOARD_WIDTH, TOP_HEIGHT);
    ctx.stroke();
    // 两侧分隔线
    ctx.beginPath();
    ctx.moveTo(SIDE_WIDTH, TOP_HEIGHT);
    ctx.lineTo(SIDE_WIDTH, BOARD_HEIGHT);
    ctx.moveTo(BOARD_WIDTH - SIDE_WIDTH, TOP_HEIGHT);
    ctx.lineTo(BOARD_WIDTH - SIDE_WIDTH, BOARD_HEIGHT);
    ctx.stroke();
    
    // 绘制格线，需要考虑顶部偏移
    ctx.lineWidth = 1;
    
    // 横线
    for (let i = 0; i < 10; i++) {
        ctx.beginPath();
        ctx.moveTo(SIDE_WIDTH + GRID_SIZE, TOP_HEIGHT + GRID_SIZE + i * GRID_SIZE);
        ctx.lineTo(SIDE_WIDTH + GRID_SIZE * 9, TOP_HEIGHT + GRID_SIZE + i * GRID_SIZE);
        ctx.stroke();
    }
    
    // 竖线
    for (let i = 0; i < 9; i++) {
        ctx.beginPath();
        ctx.moveTo(SIDE_WIDTH + GRID_SIZE + i * GRID_SIZE, TOP_HEIGHT + GRID_SIZE);
        ctx.lineTo(SIDE_WIDTH + GRID_SIZE + i * GRID_SIZE, TOP_HEIGHT + GRID_SIZE * 10);
        ctx.stroke();
    }
    
    // 添加九宫格斜线，需要考虑顶部偏移
    // 红方九宫格
    ctx.beginPath();
    ctx.moveTo(SIDE_WIDTH + GRID_SIZE * 4, TOP_HEIGHT + GRID_SIZE);
    ctx.lineTo(SIDE_WIDTH + GRID_SIZE * 6, TOP_HEIGHT + GRID_SIZE * 3);
    ctx.stroke();
    
    ctx.beginPath();
    ctx.moveTo(SIDE_WIDTH + GRID_SIZE * 6, TOP_HEIGHT + GRID_SIZE);
    ctx.lineTo(SIDE_WIDTH + GRID_SIZE * 4, TOP_HEIGHT + GRID_SIZE * 3);
    ctx.stroke();
    
    // 黑方九宫格
    ctx.beginPath();
    ctx.moveTo(SIDE_WIDTH + GRID_SIZE * 4, TOP_HEIGHT + GRID_SIZE * 8);
    ctx.lineTo(SIDE_WIDTH + GRID_SIZE * 6, TOP_HEIGHT + GRID_SIZE * 10);
    ctx.stroke();
    
    ctx.beginPath();
    ctx.moveTo(SIDE_WIDTH + GRID_SIZE * 6, TOP_HEIGHT + GRID_SIZE * 8);
    ctx.lineTo(SIDE_WIDTH + GRID_SIZE * 4, TOP_HEIGHT + GRID_SIZE * 10);
    ctx.stroke();
    
    // 绘制棋子，需要考虑顶部偏移
    for (let i = 0; i < 10; i++) {
        for (let j = 0; j < 9; j++) {
            if (gameBoard[i][j]) {
                drawPiece(j, i, gameBoard[i][j], false);
            }
        }
    }
    
    // 绘制被吃掉的棋子
    drawCapturedPieces();
    
    // 绘制游戏状态
    drawGameStatus();
}

// 绘制棋子
function drawPiece(x, y, piece, isCaptured) {
    let centerX = GRID_SIZE + x * GRID_SIZE;
    if (!isCaptured) {
        centerX += SIDE_WIDTH;
    }
    const centerY = TOP_HEIGHT + GRID_SIZE + y * GRID_SIZE;
    
    // 绘制棋子背景
    ctx.beginPath();
    ctx.arc(centerX, centerY, PIECE_RADIUS, 0, Math.PI * 2);
    ctx.fillStyle = '#fff';
    ctx.fill();
    ctx.strokeStyle = '#000';
    ctx.stroke();
    
    // 绘制棋子文字
    ctx.fillStyle = isRedPiece(piece) ? '#f00' : '#000';
    ctx.font = '24px KaiTi';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    ctx.fillText(piece, centerX, centerY);
}

// 判断是否为红方棋子
function isRedPiece(piece) {
    return '車馬相仕帥炮兵'.includes(piece);
}

// 修改点击事件监听函数
canvas.addEventListener('click', function(e) {
    const rect = canvas.getBoundingClientRect();
    const mouseX = e.clientX - rect.left;
    const mouseY = e.clientY - rect.top;
    
    // 重新开始按钮的位置和尺寸
    const buttonX = BOARD_WIDTH - 80;
    const buttonY = TOP_HEIGHT / 2;
    const buttonWidth = 100;
    const buttonHeight = 36;
    
    // 检查点击是否在按钮范围内
    if (mouseX >= buttonX - buttonWidth/2 && mouseX <= buttonX + buttonWidth/2 &&
        mouseY >= buttonY - buttonHeight/2 && mouseY <= buttonY + buttonHeight/2) {
        resetGame();
        return;
    }
    
    // 检查是否点击了棋盘，如果是则处理棋子移动
    const x = Math.round((mouseX - SIDE_WIDTH - GRID_SIZE) / GRID_SIZE);
    const y = Math.round((mouseY - TOP_HEIGHT - GRID_SIZE) / GRID_SIZE);
    
    if (x >= 0 && x < 9 && y >= 0 && y < 10) {
        handleClick(x, y);
    }
});

// 处理棋子选择和移动
function handleClick(x, y) {
    // 如果游戏已结束，不允许继续移动
    if (gameOver) {
        return;
    }

    const piece = gameBoard[y][x];
    
    if (!selectedPiece) {
        // 选择棋子
        if (piece && isCurrentPlayerPiece(piece)) {
            selectedPiece = { x, y };
            possibleMoves = calculatePossibleMoves(x, y);
            drawBoard();
            highlightSelectedPiece(x, y);
            highlightPossibleMoves();
        }
    } else {
        // 检查是否点击了可能的移动位置
        const moveIndex = possibleMoves.findIndex(move => move.x === x && move.y === y);
        if (moveIndex !== -1) {
            // 更新计时器
            const currentTime = Date.now();
            const elapsedTime = Math.floor((currentTime - moveStartTime) / 1000);
            if (currentPlayer === 'red') {
                redTotalTime += elapsedTime;
            } else {
                blackTotalTime += elapsedTime;
            }
            moveStartTime = currentTime;
            currentTurnStartTime = currentTime;  // 重置回合计时器
            
            // 如果目标位置有棋子，则添加到被吃掉的棋子列表中
            const targetPiece = gameBoard[y][x];
            if (targetPiece) {
                if (isRedPiece(targetPiece)) {
                    capturedRed.push(targetPiece);
                    // 判断是否吃掉红方帥
                    if (targetPiece === '帥') {
                        gameOver = true;
                        winner = 'black';
                    }
                } else {
                    capturedBlack.push(targetPiece);
                    // 判断是否吃掉黑方將
                    if (targetPiece === '將') {
                        gameOver = true;
                        winner = 'red';
                    }
                }
            }
            
            // 移动棋子
            gameBoard[y][x] = gameBoard[selectedPiece.y][selectedPiece.x];
            gameBoard[selectedPiece.y][selectedPiece.x] = '';
            currentPlayer = currentPlayer === 'red' ? 'black' : 'red';
        }
        selectedPiece = null;
        possibleMoves = [];
        drawBoard();
        
        // 如果游戏结束，显示获胜信息
        if (gameOver) {
            showWinner();
        }
    }
}

// 高亮选中的棋子
function highlightSelectedPiece(x, y) {
    const centerX = SIDE_WIDTH + GRID_SIZE + x * GRID_SIZE;
    const centerY = TOP_HEIGHT + GRID_SIZE + y * GRID_SIZE;
    
    ctx.beginPath();
    ctx.arc(centerX, centerY, PIECE_RADIUS + 2, 0, Math.PI * 2);
    ctx.strokeStyle = '#0f0';
    ctx.lineWidth = 2;
    ctx.stroke();
}

// 判断是否为当前玩家的棋子
function isCurrentPlayerPiece(piece) {
    return currentPlayer === 'red' ? isRedPiece(piece) : !isRedPiece(piece);
}

// 判断移动是否合法（简化版）
function isValidMove(fromX, fromY, toX, toY) {
    const piece = gameBoard[fromY][fromX];
    const targetPiece = gameBoard[toY][toX];
    
    // 不能吃自己的棋子
    if (targetPiece && isRedPiece(piece) === isRedPiece(targetPiece)) {
        return false;
    }
    
    // 检查是否在可能的移动位置列表中
    return possibleMoves.some(move => move.x === toX && move.y === toY);
}

// 重置游戏
function resetGame() {
    gameBoard = JSON.parse(JSON.stringify(initialBoard));
    selectedPiece = null;
    currentPlayer = 'red';
    capturedRed = [];
    capturedBlack = [];
    gameOver = false;
    winner = null;
    
    // 重置计时器
    moveStartTime = Date.now();
    currentTurnStartTime = Date.now();
    redTotalTime = 0;
    blackTotalTime = 0;
    drawBoard();
}

// 修改鼠标移动事件监听函数
canvas.addEventListener('mousemove', function(e) {
    // 如果已经选中棋子，就不处理鼠标移动事件
    if (selectedPiece) {
        return;
    }

    const rect = canvas.getBoundingClientRect();
    const mouseX = e.clientX - rect.left;
    const mouseY = e.clientY - rect.top;
    
    // 将鼠标位置转换为棋盘格子索引
    const boardX = mouseX - SIDE_WIDTH - GRID_SIZE;
    const boardY = mouseY - TOP_HEIGHT - GRID_SIZE;
    
    // 使用 Math.round 计算最近的格子
    const x = Math.round(boardX / GRID_SIZE);
    const y = Math.round(boardY / GRID_SIZE);
    
    // 检查位置是否在棋盘上且有棋子
    if (x >= 0 && x < 9 && y >= 0 && y < 10) {
        const piece = gameBoard[y][x];
        
        // 检查是否是当前玩家的棋子
        if (piece && isCurrentPlayerPiece(piece)) {
            // 存储悬停位置
            hoveredPosition = { x, y };
            
            // 计算可能的移动位置
            possibleMoves = calculatePossibleMoves(x, y);
            
            // 手动绘制所有内容
            drawBoard();
            
            // 高亮当前悬停的棋子
            const centerX = SIDE_WIDTH + GRID_SIZE + x * GRID_SIZE;
            const centerY = TOP_HEIGHT + GRID_SIZE + y * GRID_SIZE;
            
            ctx.beginPath();
            ctx.arc(centerX, centerY, PIECE_RADIUS + 2, 0, Math.PI * 2);
            ctx.strokeStyle = '#0f0';
            ctx.lineWidth = 2;
            ctx.stroke();
            
            // 显示可能的移动位置
            highlightPossibleMoves();
        } else {
            hoveredPosition = null;
            possibleMoves = [];
        }
    } else {
        hoveredPosition = null;
        possibleMoves = [];
    }
});

// 修改高亮显示可能移动位置的函数
function highlightPossibleMoves() {
    if (!possibleMoves || possibleMoves.length === 0) return;
    
    for (const move of possibleMoves) {
        const centerX = SIDE_WIDTH + GRID_SIZE + move.x * GRID_SIZE;
        const centerY = TOP_HEIGHT + GRID_SIZE + move.y * GRID_SIZE;
        
        ctx.beginPath();
        ctx.arc(centerX, centerY, PIECE_RADIUS / 2, 0, Math.PI * 2);
        ctx.fillStyle = 'rgba(0, 255, 0, 0.3)';
        ctx.fill();
    }
}

// 添加计算可能移动位置的函数
function calculatePossibleMoves(x, y) {
    const piece = gameBoard[y][x];
    const moves = [];
    const isRed = isRedPiece(piece);
    
    switch (piece) {
        case '兵':
        case '卒':
            // 兵卒的走法
            const direction = piece === '兵' ? 1 : -1; // 红兵向下，黑卒向上
            // 未过河
            if ((piece === '兵' && y < 5) || (piece === '卒' && y > 4)) {
                if (isValidPosition(x, y + direction)) {
                    moves.push({x: x, y: y + direction});
                }
            } else {
                // 过河后
                if (isValidPosition(x, y + direction)) {
                    moves.push({x: x, y: y + direction});
                }
                if (isValidPosition(x + 1, y)) {
                    moves.push({x: x + 1, y: y});
                }
                if (isValidPosition(x - 1, y)) {
                    moves.push({x: x - 1, y: y});
                }
            }
            break;
            
        case '車':
        case '俥':
            // 车的走法（横向和纵向）
            // 横向移动
            for (let i = 0; i < 9; i++) {
                if (i !== x) {
                    let blocked = false;
                    // 检查路径上是否有棋子
                    const start = Math.min(x, i);
                    const end = Math.max(x, i);
                    for (let j = start + 1; j < end; j++) {
                        if (gameBoard[y][j] !== '') {
                            blocked = true;
                            break;
                        }
                    }
                    if (!blocked) {
                        moves.push({x: i, y: y});
                    }
                }
            }
            // 纵向移动
            for (let i = 0; i < 10; i++) {
                if (i !== y) {
                    let blocked = false;
                    // 检查路径上是否有棋子
                    const start = Math.min(y, i);
                    const end = Math.max(y, i);
                    for (let j = start + 1; j < end; j++) {
                        if (gameBoard[j][x] !== '') {
                            blocked = true;
                            break;
                        }
                    }
                    if (!blocked) {
                        moves.push({x: x, y: i});
                    }
                }
            }
            break;
            
        case '馬':
        case '傌':
            // 马的八个可能移动位置
            const horseOffsets = [
                {x: -2, y: -1}, {x: -2, y: 1},
                {x: 2, y: -1}, {x: 2, y: 1},
                {x: -1, y: -2}, {x: -1, y: 2},
                {x: 1, y: -2}, {x: 1, y: 2}
            ];
            
            for (const offset of horseOffsets) {
                const newX = x + offset.x;
                const newY = y + offset.y;
                
                // 检查马腿
                if (isValidPosition(newX, newY)) {
                    let blocked = false;
                    if (Math.abs(offset.x) === 2) {
                        blocked = gameBoard[y][x + offset.x/2] !== '';
                    } else {
                        blocked = gameBoard[y + offset.y/2][x] !== '';
                    }
                    
                    if (!blocked) {
                        moves.push({x: newX, y: newY});
                    }
                }
            }
            break;

        case '相':
        case '象':
            // 象的四个可能移动位置
            const elephantOffsets = [
                {x: -2, y: -2}, {x: 2, y: -2},
                {x: -2, y: 2}, {x: 2, y: 2}
            ];
            
            for (const offset of elephantOffsets) {
                const newX = x + offset.x;
                const newY = y + offset.y;
                
                // 检查象眼和不能过河
                if (isValidPosition(newX, newY) && 
                    ((isRed && newY <= 4) || (!isRed && newY >= 5))) {
                    // 检查象眼
                    const eyeX = x + offset.x/2;
                    const eyeY = y + offset.y/2;
                    if (gameBoard[eyeY][eyeX] === '') {
                        moves.push({x: newX, y: newY});
                    }
                }
            }
            break;

        case '仕':
        case '士':
            // 士的四个可能移动位置
            const advisorOffsets = [
                {x: -1, y: -1}, {x: 1, y: -1},
                {x: -1, y: 1}, {x: 1, y: 1}
            ];
            
            for (const offset of advisorOffsets) {
                const newX = x + offset.x;
                const newY = y + offset.y;
                
                // 检查是否在九宫格内
                if (newX >= 3 && newX <= 5 && 
                    ((isRed && newY >= 0 && newY <= 2) || 
                     (!isRed && newY >= 7 && newY <= 9))) {
                    moves.push({x: newX, y: newY});
                }
            }
            break;

        case '帥':
        case '將':
            // 将/帅的四个可能移动位置
            const kingOffsets = [
                {x: 0, y: -1}, {x: 0, y: 1},
                {x: -1, y: 0}, {x: 1, y: 0}
            ];
            
            for (const offset of kingOffsets) {
                const newX = x + offset.x;
                const newY = y + offset.y;
                
                // 检查是否在九宫格内
                if (newX >= 3 && newX <= 5 && 
                    ((isRed && newY >= 0 && newY <= 2) || 
                     (!isRed && newY >= 7 && newY <= 9))) {
                    moves.push({x: newX, y: newY});
                }
            }
            
            // 检查是否可以直接吃对方的将/帥
            // 在同一列上寻找对方的将/帥
            let foundOpponentKing = false;
            let hasBlockingPiece = false;
            
            // 确定搜索的起点和终点
            const startY = isRed ? y + 1 : 0;
            const endY = isRed ? 9 : y - 1;
            
            // 在同一列上搜索
            for (let i = startY; i <= endY; i++) {
                if (gameBoard[i][x] !== '') {
                    if ((isRed && gameBoard[i][x] === '將') || 
                        (!isRed && gameBoard[i][x] === '帥')) {
                        foundOpponentKing = true;
                    } else {
                        hasBlockingPiece = true;
                    }
                    break;
                }
            }
            
            // 如果找到对方的将/帥且中间没有其他棋子，添加这个位置
            if (foundOpponentKing && !hasBlockingPiece) {
                const targetY = isRed ? 9 : 0;
                moves.push({x: x, y: targetY});
            }
            break;

        case '炮':
        case '砲':
            // 炮的走法
            // 横向移动
            for (let i = 0; i < 9; i++) {
                if (i !== x) {
                    let jumpCount = 0;
                    // 计算路径上的棋子数
                    const start = Math.min(x, i);
                    const end = Math.max(x, i);
                    for (let j = start + 1; j < end; j++) {
                        if (gameBoard[y][j] !== '') {
                            jumpCount++;
                        }
                    }
                    // 根据跳子数决定是否可以移动
                    if ((jumpCount === 0 && gameBoard[y][i] === '') || 
                        (jumpCount === 1 && gameBoard[y][i] !== '')) {
                        moves.push({x: i, y: y});
                    }
                }
            }
            // 纵向移动
            for (let i = 0; i < 10; i++) {
                if (i !== y) {
                    let jumpCount = 0;
                    // 计算路径上的棋子数
                    const start = Math.min(y, i);
                    const end = Math.max(y, i);
                    for (let j = start + 1; j < end; j++) {
                        if (gameBoard[j][x] !== '') {
                            jumpCount++;
                        }
                    }
                    // 根据跳子数决定是否可以移动
                    if ((jumpCount === 0 && gameBoard[i][x] === '') || 
                        (jumpCount === 1 && gameBoard[i][x] !== '')) {
                        moves.push({x: x, y: i});
                    }
                }
            }
            break;
    }
    
    // 过滤掉不合法的移动（比如吃自己的子）
    return moves.filter(move => {
        const targetPiece = gameBoard[move.y][move.x];
        return !targetPiece || isRedPiece(targetPiece) !== isRed;
    });
}

// 添加位置验证函数
function isValidPosition(x, y) {
    return x >= 0 && x < 9 && y >= 0 && y < 10;
}

// 修改绘制被吃掉棋子的函数
function drawCapturedPieces() {
    const PIECES_PER_COLUMN = 8; // 每列最多显示8个棋子
    const COLUMN_SPACING = SIDE_WIDTH / 3; // 修改列间距为区域宽度的1/3
    const MARGIN = 20; // 添加边距，避免棋子贴近分隔线
    
    // 绘制左侧红色棋子
    capturedRed.forEach((piece, index) => {
        const column = Math.floor(index / PIECES_PER_COLUMN);
        const row = index % PIECES_PER_COLUMN;
        
        const centerX = MARGIN + COLUMN_SPACING + column * COLUMN_SPACING;
        const centerY = TOP_HEIGHT + GRID_SIZE + row * (PIECE_RADIUS * 2 + 5);
        
        // 绘制棋子背景
        ctx.beginPath();
        ctx.arc(centerX, centerY, PIECE_RADIUS * 0.8, 0, Math.PI * 2);
        ctx.fillStyle = '#fff';
        ctx.fill();
        ctx.strokeStyle = '#000';
        ctx.stroke();
        
        // 绘制棋子文字
        ctx.fillStyle = '#f00';
        ctx.font = '20px KaiTi';
        ctx.textAlign = 'center';
        ctx.textBaseline = 'middle';
        ctx.fillText(piece, centerX, centerY);
    });
    
    // 绘制右侧黑色棋子
    capturedBlack.forEach((piece, index) => {
        const column = Math.floor(index / PIECES_PER_COLUMN);
        const row = index % PIECES_PER_COLUMN;
        
        const centerX = BOARD_WIDTH - MARGIN - COLUMN_SPACING - column * COLUMN_SPACING;
        const centerY = TOP_HEIGHT + GRID_SIZE + row * (PIECE_RADIUS * 2 + 5);
        
        // 绘制棋子背景
        ctx.beginPath();
        ctx.arc(centerX, centerY, PIECE_RADIUS * 0.8, 0, Math.PI * 2);
        ctx.fillStyle = '#fff';
        ctx.fill();
        ctx.strokeStyle = '#000';
        ctx.stroke();
        
        // 绘制棋子文字
        ctx.fillStyle = '#000';
        ctx.font = '20px KaiTi';
        ctx.textAlign = 'center';
        ctx.textBaseline = 'middle';
        ctx.fillText(piece, centerX, centerY);
    });
}

// 添加显示获胜信息的函数
function showWinner() {
    ctx.fillStyle = 'rgba(0, 0, 0, 0.5)';
    ctx.fillRect(0, 0, BOARD_WIDTH, BOARD_HEIGHT);
    
    ctx.fillStyle = '#fff';
    ctx.font = 'bold 48px KaiTi';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    const text = winner === 'red' ? '红方胜利！' : '黑方胜利！';
    ctx.fillText(text, BOARD_WIDTH / 2, BOARD_HEIGHT / 2);
}

// 修改 drawGameStatus 函数
function drawGameStatus() {
    const currentTime = Date.now();
    const elapsedTime = Math.floor((currentTime - moveStartTime) / 1000);
    const turnTime = Math.floor((currentTime - currentTurnStartTime) / 1000);
    
    // 显示当前回合方
    ctx.font = 'bold 24px KaiTi';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    ctx.fillStyle = currentPlayer === 'red' ? '#f00' : '#000';
    const playerText = currentPlayer === 'red' ? '红方回合' : '黑方回合';
    ctx.fillText(playerText, BOARD_WIDTH / 2, TOP_HEIGHT / 3);
    
    // 显示双方总用时
    ctx.font = '20px KaiTi';
    // 红方总用时
    ctx.fillStyle = '#f00';
    const redTimeText = `总用时: ${formatTime(redTotalTime + (currentPlayer === 'red' ? elapsedTime : 0))}`;
    ctx.fillText(redTimeText, SIDE_WIDTH + GRID_SIZE * 2, TOP_HEIGHT / 3);
    
    // 黑方总用时
    ctx.fillStyle = '#000';
    const blackTimeText = `总用时: ${formatTime(blackTotalTime + (currentPlayer === 'black' ? elapsedTime : 0))}`;
    ctx.fillText(blackTimeText, BOARD_WIDTH - SIDE_WIDTH - GRID_SIZE * 2, TOP_HEIGHT / 3);
    
    // 显示当前回合用时
    ctx.font = '18px KaiTi';
    ctx.fillStyle = currentPlayer === 'red' ? '#f00' : '#000';
    const turnTimeText = `本回合: ${formatTime(turnTime)}`;
    ctx.fillText(turnTimeText, BOARD_WIDTH / 2, TOP_HEIGHT * 2/3);
    
    // 绘制重新开始按钮
    const buttonX = BOARD_WIDTH - 80;
    const buttonY = TOP_HEIGHT / 2;
    const buttonWidth = 100;
    const buttonHeight = 36;
    
    // 绘制按钮背景 - 使用与棋盘相协调的颜色
    ctx.fillStyle = '#d28b48';  // 使用暖色调的棕色，与棋盘风格相匹配
    ctx.beginPath();
    ctx.rect(buttonX - buttonWidth/2, buttonY - buttonHeight/2, buttonWidth, buttonHeight);
    ctx.fill();
    
    // 绘制按钮边框
    ctx.strokeStyle = '#8c5e2a';  // 深棕色边框
    ctx.lineWidth = 2;
    ctx.stroke();
    
    // 绘制按钮文字
    ctx.fillStyle = '#fff';  // 保持白色文字，提供良好的对比度
    ctx.font = 'bold 16px KaiTi';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    ctx.fillText('重新开始', buttonX, buttonY);
}

// 添加时间格式化函数
function formatTime(seconds) {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
}

// 修改 updateTimer 函数
function updateTimer() {
    if (!gameOver && !hoveredPosition) {
        drawBoard();
    }
    requestAnimationFrame(updateTimer);
}

// 启动计时器
updateTimer();

// 初始化游戏
drawBoard(); 