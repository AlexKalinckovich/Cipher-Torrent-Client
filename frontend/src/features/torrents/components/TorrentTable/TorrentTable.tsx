import React, { useCallback } from 'react';
import { Table } from 'antd';
import type { HTMLAttributes, MouseEvent } from 'react';
import type { Torrent } from '@/types/model/models.ts';
import type { TorrentTableProps } from '@/features/torrents/types/dashboardTypes';
import styles from './TorrentTable.module.css';

const checkIgnoredClick = (target: HTMLElement): boolean => {
    const isButton: boolean = target.closest('button') !== null;
    const isSvg: boolean = target.closest('svg') !== null;
    return isButton || isSvg;
};

const getRowClassName = (): string => {
    return styles.clickableRow;
};

export const TorrentTable: React.FC<TorrentTableProps> = ({ data, loading, columns, onRowClick }) => {
    const paginationConfig: { pageSize: number } = {
        pageSize: 10
    };

    const handleRow = useCallback((record: Torrent): HTMLAttributes<HTMLElement> => {
        const handleClick = (e: MouseEvent<HTMLElement>): void => {
            const target: HTMLElement = e.target as HTMLElement;

            if (checkIgnoredClick(target)) {
                return;
            }

            onRowClick(record);
        };

        return {
            onClick: handleClick
        };
    }, [onRowClick]);

    return (
        <div className={styles.tableWrapper}>
            <Table<Torrent>
                dataSource={data}
                columns={columns}
                rowKey="info_hash"
                loading={loading}
                pagination={paginationConfig}
                onRow={handleRow}
                rowClassName={getRowClassName}
            />
        </div>
    );
};