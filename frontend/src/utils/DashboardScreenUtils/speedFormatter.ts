const UNITS: string[] = ['B/s', 'KB/s', 'MB/s', 'GB/s', 'TB/s'];

const calculateIndex = (value: number): number => {
    return Math.floor(Math.log(value) / Math.log(1024));
};

const getSafeIndex = (value: number): number => {
    if (value === 0) {
        return 0;
    }
    return calculateIndex(value);
};

export const formatSpeed = (bps: number | undefined): string => {
    const safeBps: number = bps ?? 0;
    const index: number = getSafeIndex(safeBps);
    const size: number = safeBps / Math.pow(1024, index);
    return `${size.toFixed(1)} ${UNITS[index]}`;
};