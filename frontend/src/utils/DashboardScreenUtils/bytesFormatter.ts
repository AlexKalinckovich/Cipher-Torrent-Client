
const BYTE_SIZES: string[] = ['B', 'KB', 'MB', 'GB', 'TB'];
const KILOBYTE: number = 1024;
const FLOAT_NUMBER_COUNT: number = 2;

export const formatBytes = (bytes: number): string => {
    if (bytes === 0) return '0 B';
    return buildByteString(bytes, getByteIndex(bytes));
};

const getByteIndex = (bytes: number): number => {
    return Math.floor(Math.log(bytes) / Math.log(KILOBYTE));
};

const buildByteString = (bytes: number, index: number): string => {
    const value = parseFloat((bytes / Math.pow(KILOBYTE, index)).toFixed(FLOAT_NUMBER_COUNT));
    return `${value} ${BYTE_SIZES[index]}`;
};

