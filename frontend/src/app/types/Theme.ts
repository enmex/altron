export type Theme = {
    name: string;
    primary: string;
    text: string;
    secondary: string;
    tertiary: string;
    accents: {
        positive: string;
        negative: string;
        neutral: string;
        contrast: string;
    };
    webpPath?: string;
    greetingsSoundPath?: string;
    request: string;
    response: string;
}