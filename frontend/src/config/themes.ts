import { Theme } from "../app/types/Theme";

export const THEMES: Theme[] = [
    {
        name: "Dark",
        primary: "#1E1E1E",
        text: "#D4D4D4",
        secondary: "#2A2A2A",
        tertiary: "#3C3C3C",
        accents: {
            positive: "#23D18B",
            negative: "#F14C4C",
            neutral: "#D4D4D4",
            contrast: "#FFFFFF"
        },
        webpPath: "/assets/dark.webp",
        greetingsSoundPath: "/assets/sounds/greetings_default.wav",
        request: "#2B5876",
        response: "#399247"
    },
    {
        name: "Light",
        primary: "#f8f8f8",            
        text: "#333333",             
        secondary: "#ffffff",        
        tertiary: "#f0f0f0",         
        accents: {                   
            positive: "#64a844",     
            negative: "#e74c3c",     
            neutral: "#3498db",      
            contrast: "#ff9900",     
        },
        greetingsSoundPath: "/assets/sounds/greetings_default.wav",
        request: "#FFA726", 
        response: "#1E88E5"
    },
    {
        name: "Retrowave",
        primary: "#0E0E10",      // Основной фон
        text: "#FFFFFF",          // Цвет текста
        secondary: "#1A1A1D",    // Вторичный фон (для диалогового окна)
        tertiary: "#7e0089",      // Третичный фон (для контраста)
        accents: {
            positive: "#009de1",    // Цвет для позитивных элементов
            negative: "#8308a1",    // Цвет для негативных элементов
            neutral: "#e10099",     // Нейтральный цвет
            contrast: "#e0b82f"     // Цвет для выделения и контраста
        },
        webpPath: "/assets/retrowave.webp",  // Путь к фоновой гифке (если есть)
        greetingsSoundPath: "/assets/sounds/greetings_retrowave.wav",
        request: "#32023a",  // Розовый цвет диалогового окна для запросов
        response: "#02133a"  // Желтый цвет диалогового окна для ответов

    },
    {
        name: "Sith Lord",
        primary: "#ffffff",
        text: "#000000",
        secondary: "#dfdfdf",
        tertiary: "#ff0000",
        accents: {
            positive: "#ff0000",
            negative: "#394A59",
            neutral: "#2d2b3c",
            contrast: "#2d2c2c" 
        },
        webpPath: "/assets/starwars.webp",
        greetingsSoundPath: "/assets/sounds/greetings_starwars.wav",
        request: "#dfdfdf",
        response: "#dfdfdf"
    },
    {
        name: "High Contrast",
        text: "#FFFFFF",
        primary: "#000000",
        secondary: "#1A1A1A",
        tertiary: "#333333", 
        accents: {
            positive: "#00FF00",
            negative: "#FF0000",
            neutral: "#FFFF00", 
            contrast: "#FFA500",
        },
        greetingsSoundPath: "/assets/sounds/greetings_default.wav",
        request: "#27AE60",
        response: "#FFB450"
    },
    {
        name: "Matrix",
        primary: "#000000",
        text: "#00FF00",
        secondary: "#080808",
        tertiary: "#0F0F0F",
        accents: {
            positive: "#00FF00",
            negative: "#FF0000",
            neutral: "#808080",
            contrast: "#00FFFF"
        },
        webpPath: "/assets/matrix.webp",
        greetingsSoundPath: "/assets/sounds/greetings_matrix.wav",
        request: "#080808",
        response: "#111111"
    },
    {
        name: "Windows95",
        primary: "#008080", // Teal
        text: "#000000", // Black
        secondary: "#C0C0C0", // Silver
        tertiary: "#FFFFFF", // White
        accents: {
            positive: "#00FF00", // Lime
            negative: "#FF0000", // Red
            neutral: "#808080", // Gray
            contrast: "#FFFF00", // Yellow
        },
        webpPath: "/assets/windows95.webp",
        greetingsSoundPath: "/assets/sounds/greetings_windows95.wav",
        request: "#F0F0F0", 
        response: "#F0F0F0",
    },
    {
        name: "Vault-Tech",
        primary: "#152340",     // Темно-синий цвет для фона
        text: "#FFFFFF",        // Белый цвет для текста
        secondary: "#2B5490",   // Голубой цвет для диалоговых панелей "запроса"
        tertiary: "#1E3C6E",    // Ещё более тёмный голубой цвет для других элементов интерфейса
        accents: {
            positive: "#6ED7A2", // Зелёный цвет для положительных акцентов
            negative: "#E4572E", // Красный цвет для отрицательных акцентов
            neutral: "#FFD166",  // Жёлтый цвет для нейтральных акцентов
            contrast: "#F7A8B8", // Розовый цвет для контраста и выделения
        },
        webpPath: "/assets/fallout.webp",
        greetingsSoundPath: "/assets/sounds/greetings_default.wav",
        request: "#DF6B7B",    // Светло-красный цвет для диалоговых панелей "запроса"
        response: "#8FB1CC",   // Светло-голубой цвет для диалоговых панелей "ответа"
    },
    {
        name: "Barbie Glamour",
        primary: "#F9B0C9",     
        text: "#000000",        
        secondary: "#FDD9EB",   
        tertiary: "#F6A8C4",    
        accents: {
            positive: "#FF69B4", 
            negative: "#9C27B0", 
            neutral: "#D16002",  
            contrast: "#03A9F4", 
        },
        greetingsSoundPath: "/assets/sounds/greetings_default.wav",
        request: "#FF69B4",    
        response: "#9C27B0",   
    },
    {
        name: "Cyberpunk",
        primary: "#080216",    // Основной фон
        text: "#1eddff",        // Цвет текста
        secondary: "#290a0c",  // Вторичный фон (для диалогового окна)
        tertiary: "#c21818",    // Третичный фон (для контраста)
        accents: {
            positive: "#3bb04b",  // Цвет для позитивных элементов
            negative: "#db3a55",  // Цвет для негативных элементов
            neutral: "#d4eb10",   // Нейтральный цвет
            contrast: "#d4eb10"   // Цвет для выделения и контраста
        },
        webpPath: "/assets/cyberpunk.webp",  // Путь к фоновой гифке (если есть)
        greetingsSoundPath: "/assets/sounds/greetings_cyberpunk.wav",
        request: "#060414",  // Красный цвет диалогового окна для запросов
        response: "#083639"  // Желтый цвет диалогового окна для ответов
    }
];