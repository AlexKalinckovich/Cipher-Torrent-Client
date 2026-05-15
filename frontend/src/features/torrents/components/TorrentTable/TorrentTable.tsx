import React from 'react';
import { Table } from 'antd';
import type { Torrent } from '@/types/model/models.ts';
import type { TorrentTableProps } from '@/features/torrents/types/dashboardTypes';
import styles from './TorrentTable.module.css';

export const TorrentTable: React.FC<TorrentTableProps> = ({ data, loading, columns }) => {
    const paginationConfig: { pageSize: number } = {
        pageSize: 10
    };
    return (
        <div className={styles.tableWrapper}>
            <Table<Torrent>
                dataSource={data}
                columns={columns}
                rowKey="info_hash"
                loading={loading}
                pagination={paginationConfig}
            />
        </div>
    );
};
