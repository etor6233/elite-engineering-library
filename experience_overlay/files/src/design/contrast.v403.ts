// AUTHORED verification glue implementing W3C relative luminance.
export function contrast(foreground: string, background: string): number {
    const luminance = (hex: string) => {
        if (!/^#[a-f0-9]{6}$/i.test(hex))
            throw new Error("Six-digit sRGB required");
        const rgb = [1, 3, 5].map(i => parseInt(hex.slice(i, i + 2), 16) / 255).map(c => c <= .04045 ? c / 12.92 : ((c + .055) / 1.055) ** 2.4);
        return rgb[0]! * 0.2126 + rgb[1]! * 0.7152 + rgb[2]! * 0.0722;
    };
    const a = luminance(foreground), b = luminance(background);
    return (Math.max(a, b) + .05) / (Math.min(a, b) + .05);
}
