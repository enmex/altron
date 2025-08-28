import {Buffer} from 'buffer';

const dictionary = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";

export const mergeArrays = <T extends {}>(a: T[], b: T[]): T[] => {
    const mergedArray: T[] = [];

    for (let i = 0; i < Math.max(a.length, b.length); i++) {
        if (i < a.length) {
            mergedArray.push(a[i]);
        }
        if (i < b.length) {
            mergedArray.push(b[i]);
        }
    }

    return mergedArray;
}

export const randomKey = (): string => {
    let res = "";
    for (let i = 0; i < 20; i++) {
        res += dictionary[Math.floor(Math.random() * dictionary.length)];
    }
    return res;
}

export const darkenColor = (color: string, factor: number): string => {
  if (factor < 0 || factor > 1) {
    return color;
  }

  const r = parseInt(color.slice(1, 3), 16);
  const g = parseInt(color.slice(3, 5), 16);
  const b = parseInt(color.slice(5, 7), 16);

  let darkenedR = Math.round(r * factor);
  let darkenedG = Math.round(g * factor);
  let darkenedB = Math.round(b * factor);

  const darkenedHex = `#${darkenedR.toString(16).padStart(2, '0')}${darkenedG.toString(16).padStart(2, '0')}${darkenedB.toString(16).padStart(2, '0')}`;

  return darkenedHex;
}

export const negativeColor = (color: string): string => {
    const r = parseInt(color.slice(1, 3), 16);
    const g = parseInt(color.slice(3, 5), 16);
    const b = parseInt(color.slice(5, 7), 16);

    return `#${(255 - r).toString(16)}${(255 - g).toString(16)}${(255 - b).toString(16)}`
}

export const blendColors = (color1: string, color2: string, ratio: number): string => {
    const r1 = parseInt(color1.substr(1, 2), 16);
    const g1 = parseInt(color1.substr(3, 2), 16);
    const b1 = parseInt(color1.substr(5, 2), 16);
    
    const r2 = parseInt(color2.substr(1, 2), 16);
    const g2 = parseInt(color2.substr(3, 2), 16);
    const b2 = parseInt(color2.substr(5, 2), 16);
    
    const r = Math.round(r1 * (1 - ratio) + r2 * ratio);
    const g = Math.round(g1 * (1 - ratio) + g2 * ratio);
    const b = Math.round(b1 * (1 - ratio) + b2 * ratio);
    
    const blendedColor = "#" + ((r << 16) | (g << 8) | b).toString(16).padStart(6, "0");

    return blendedColor;
}

export const b64ToString = (b64String: string): string => {
    return Buffer.from(b64String, 'base64').toString('utf-8');
}

export const stringToB64 = (str: string): string => {
    return Buffer.from(str, 'binary').toString('base64');
}

export const textToHexDump = (text: string) => {
    const length = text.length;
    const bytes = new Uint8Array(length);
    for (let i = 0; i < length; i++) {
        bytes[i] = text.charCodeAt(i);
    }
    const decoder = new TextDecoder();
    const buffer = decoder.decode(bytes);
    const blockSize = 16;
    let lines = [];
    const hex = '0123456789ABCDEF';
    for (let b = 0; b < buffer.length; b += blockSize) {
        const block = buffer.slice(b, Math.min(b + blockSize, buffer.length));
        const addr = ('0000000000' + b.toString(10)).slice(-10);
        let codes = block.split('').map(ch => {
            const code = ch.charCodeAt(0);
            return ' ' + hex[(0xF0 & code) >> 4] + hex[0x0F & code];
        }).join('');
        codes += ' ..'.repeat(blockSize - block.length);
        // eslint-disable-next-line no-control-regex
        let chars = block.replace(/[\x00-\x1F]/g, '.');
        chars += ' '.repeat(blockSize - block.length);
        lines.push(addr + ':' + codes + ' |' + chars + '|');
    }
    return lines.join('\n');
}

export const b64ToPythonBytes = (b64: string): string => {
    let decodedString = atob(b64);
    let bytes = new Uint8Array(decodedString.length);
    for (let i = 0; i < decodedString.length; i++) {
        bytes[i] = decodedString.charCodeAt(i);
    }

    let hexBytes = Array.from(bytes).map(byte => '\\x' + byte.toString(16).padStart(2, '0'));

    return "b'" + hexBytes.join('') + "'";
}

export const getNearestHour = (minutesOffset?: number): Date => {
    const now = new Date();
    now.setHours(now.getHours() + 1, minutesOffset ?? 0, 0, 0);
    return now;
}

export const stringifyMap = (map: Map<string, any>): string => {
    const jsonObject: {[key: string]: any} = {};

    map.forEach((value, key) => {
        jsonObject[key] = value;
    });

    return JSON.stringify(jsonObject);
}

export const convertBytesToLargerUnit = (bytes: number): string => {
    if (bytes >= 1073741824) {
        return `${(bytes / 1073741824).toFixed(2)} GB`;
    }
    if (bytes >= 1048576) {
        return `${(bytes / 1048576).toFixed(2)} MB`;
    }
    return `${(bytes / 1024).toFixed(2)} KB`;
}