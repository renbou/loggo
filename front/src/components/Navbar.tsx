import React, { useRef, useCallback, useState } from "react";
import DateRangePicker from "flowbite-datepicker/DateRangePicker";

type Props = {
  onSearch: (search: string, from: Date, to: Date) => void;
};

function Navbar(props: Props) {
  const dateRangePicker = useRef<DateRangePicker>();
  const [search, setSearch] = useState("");

  // Reload the DateRangePicker once a re-render occurs
  const setDateRangePickerRef = useCallback((node: HTMLDivElement) => {
    dateRangePicker.current?.destroy();
    if (node) {
      dateRangePicker.current = new DateRangePicker(node, {
        todayBtn: true,
        todayBtnMode: 1,
        clearBtn: true,
      });
    } else {
      dateRangePicker.current = undefined;
    }

    dateRangePicker.current?.getDates;
  }, []);

  // Perform search whenever "Search" is clicked using the current search value and chosen dates
  function performSearch(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (!dateRangePicker.current) {
      return;
    }

    let [from, to] = dateRangePicker.current.getDates();
    if (from === undefined) {
      let yesterday = new Date();
      yesterday.setDate(yesterday.getDate() - 1);
      from = yesterday;
    }

    // TODO: remove once realtime log streaming is supported
    if (to === undefined) {
      let tomorrow = new Date();
      tomorrow.setDate(tomorrow.getDate() + 1);
      to = new Date();
    }

    dateRangePicker.current.setDates(from, to);
    props.onSearch(search, from, to);
  }

  return (
    <div className="bg-slate-100 p-2">
      <div className="mx-auto max-w-6xl flex gap-2">
        <div className="relative max-w-sm">
          <div className="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
            <svg
              aria-hidden="true"
              className="w-5 h-5 text-gray-500 dark:text-gray-400"
              fill="currentColor"
              viewBox="0 0 20 20"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                fillRule="evenodd"
                d="M6 2a1 1 0 00-1 1v1H4a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-1V3a1 1 0 10-2 0v1H7V3a1 1 0 00-1-1zm0 5a1 1 0 000 2h8a1 1 0 100-2H6z"
                clipRule="evenodd"
              ></path>
            </svg>
          </div>
          <div className="flex" ref={setDateRangePickerRef}>
            <input
              data-datepicker
              type="text"
              className={
                "bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg rounded-r-none focus:ring-blue-500 focus:border-blue-500 block w-36 pl-10 p-4 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
              }
              placeholder="From"
            />
            <input
              data-datepicker
              type="text"
              className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg rounded-l-none border-l-0 focus:ring-blue-500 focus:border-blue-500 block w-36 p-4 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
              placeholder="To"
            />
          </div>
        </div>
        <form className="relative flex-grow" onSubmit={performSearch}>
          <div className="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
            <svg
              aria-hidden="true"
              className="w-5 h-5 text-gray-500 dark:text-gray-400"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
              ></path>
            </svg>
          </div>
          <input
            type="search"
            id="default-search"
            className="block w-full p-4 pl-10 text-sm text-gray-900 border border-gray-300 rounded-lg bg-gray-50 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
            placeholder="Search through logs..."
            onChange={(e) => {
              setSearch(e.target.value);
            }}
            value={search}
          />
          <button
            type="submit"
            className="text-white absolute right-2.5 bottom-2.5 bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-4 py-2 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
          >
            Search
          </button>
        </form>
      </div>
    </div>
  );
}

export default Navbar;
