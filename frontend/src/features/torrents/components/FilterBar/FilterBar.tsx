import React, { useCallback } from 'react';
import { Input } from 'antd';
import { SearchOutlined } from '@ant-design/icons';
import type { FilterBarProps, FilterButtonProps } from '@/features/torrents/types/dashboardTypes';
import styles from './FilterBar.module.css';

const FilterButton: React.FC<FilterButtonProps> = ({ label, isActive, onClick, baseClassName, activeClassName }) => {
    let className: string = baseClassName;
    if (isActive) {
        className = activeClassName;
    }
    return (
        <button type="button" className={className} onClick={onClick}>
            {label}
        </button>
    );
};

export const FilterBar: React.FC<FilterBarProps> = ({ currentFilter, searchText, onFilterChange, onSearchChange, baseClassName, activeClassName }) => {
    const handleFilterAll = useCallback(() => onFilterChange(null), [onFilterChange]);
    const handleFilterDownloading = useCallback(() => onFilterChange('downloading'), [onFilterChange]);
    const handleFilterSeeding = useCallback(() => onFilterChange('seeding'), [onFilterChange]);
    const handleFilterPaused = useCallback(() => onFilterChange('paused'), [onFilterChange]);
    const handleFilterError = useCallback(() => onFilterChange('error'), [onFilterChange]);

    const handleSearchInput = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
        onSearchChange(e.target.value);
    }, [onSearchChange]);

    return (
        <div className={styles.filtersBar}>
            <div className={styles.filterButtons}>
                <FilterButton label="All" isActive={currentFilter === null} onClick={handleFilterAll} baseClassName={baseClassName} activeClassName={activeClassName} />
                <FilterButton label="Downloading" isActive={currentFilter === 'downloading'} onClick={handleFilterDownloading} baseClassName={baseClassName} activeClassName={activeClassName} />
                <FilterButton label="Seeding" isActive={currentFilter === 'seeding'} onClick={handleFilterSeeding} baseClassName={baseClassName} activeClassName={activeClassName} />
                <FilterButton label="Paused" isActive={currentFilter === 'paused'} onClick={handleFilterPaused} baseClassName={baseClassName} activeClassName={activeClassName} />
                <FilterButton label="Error" isActive={currentFilter === 'error'} onClick={handleFilterError} baseClassName={baseClassName} activeClassName={activeClassName} />
            </div>
            <Input
                placeholder="Search torrents..."
                prefix={<SearchOutlined style={{ color: '#8b5cf6' }} />}
                value={searchText}
                onChange={handleSearchInput}
                className={styles.searchInput}
            />
        </div>
    );
};
