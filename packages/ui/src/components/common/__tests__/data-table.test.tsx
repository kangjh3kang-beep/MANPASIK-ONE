import React from 'react';
import { describe, it, expect, vi } from 'vitest';
import { render, fireEvent } from '@testing-library/react';
import { DataTable } from '../data-table';

describe('DataTable Component', () => {
  const columns = [
    { key: 'id', label: 'ID' },
    { key: 'name', label: '차트명' },
    { key: 'status', label: '상태', render: (val: string) => <span data-testid="custom-status">{val.toUpperCase()}</span> }
  ];

  const data = [
    { id: 1, name: 'Audit Log 1', status: 'completed' },
    { id: 2, name: 'Audit Log 2', status: 'pending' },
    { id: 3, name: 'System Check', status: 'failed' },
    { id: 4, name: 'User Login', status: 'completed' },
    { id: 5, name: 'Data Export', status: 'completed' },
    { id: 6, name: 'External API', status: 'pending' },
  ];

  it('컬럼 헤더와 데이터가 정상적으로 렌더링되어야 한다', () => {
    const { getByText, getAllByTestId } = render(
      <DataTable data={data} columns={columns} itemsPerPage={5} />
    );
    
    expect(getByText('ID')).toBeTruthy();
    expect(getByText('차트명')).toBeTruthy();
    expect(getByText('Audit Log 1')).toBeTruthy();
    
    // 커스텀 렌더러 확인 (uppercase 변환)
    const statusElements = getAllByTestId('custom-status');
    expect(statusElements[0].textContent).toBe('COMPLETED');
  });

  it('검색 필터가 작동해야 한다', () => {
    const { getByPlaceholderText, getByText, queryByText } = render(
      <DataTable data={data} columns={columns} itemsPerPage={10} />
    );
    
    const searchInput = getByPlaceholderText('데이터 검색...');
    fireEvent.change(searchInput, { target: { value: 'System' } });
    
    expect(getByText('System Check')).toBeTruthy();
    expect(queryByText('Audit Log 1')).toBeNull();
  });

  it('페이지네이션이 작동해야 한다', () => {
    // itemsPerPage를 2로 설정하여 3페이지 생성
    const { getByText, queryByText, getByRole } = render(
      <DataTable data={data} columns={columns} itemsPerPage={2} />
    );
    
    // 1페이지 데이터 존재 확인
    expect(getByText('Audit Log 1')).toBeTruthy();
    expect(queryByText('System Check')).toBeNull(); // 3번 데이터라 2페이지 이후에 있어야 함
    
    // 2번 버튼 클릭 (2페이지 이동)
    const page2Button = getByRole('button', { name: '2' });
    fireEvent.click(page2Button);
    
    expect(getByText('System Check')).toBeTruthy();
    expect(queryByText('Audit Log 1')).toBeNull();
  });

  it('CSV 내보내기 버튼이 클릭되어야 한다', () => {
    const exportSpy = vi.fn();
    const { getByText } = render(
      <DataTable data={data} columns={columns} onExport={exportSpy} />
    );
    
    const exportButton = getByText('내보내기 (CSV)');
    fireEvent.click(exportButton);
    
    expect(exportSpy).toHaveBeenCalled();
  });
});
