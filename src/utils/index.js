export function getTokenFromLocalStorage() {
	const token = localStorage.getItem('token');
	return token;
}

export function removeTokenFromLocalStorage() {
	localStorage.removeItem('token');
	localStorage.removeItem('refresh_token');
	localStorage.removeItem('role');
	localStorage.removeItem('id');
}

export function toCommaSeparatedValues(values) {
	return values.join(',');
}

export function getRoleFromLocalStorage() {
	const role = localStorage.getItem('role');
	return role;
}

export function getRange(start, end) {
	const length = end - start + 1;
	return Array.from({ length }, (_, index) => index + start);
}

const DATES = [
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22,
	23, 24, 25, 26, 27, 28, 29, 30, 31,
];

export function createOptionsOfDates() {
	return DATES.map(date => {
		if (date < 10) {
			// if date is less than 10, add 0 before the date
			return { label: '0' + date, value: '0' + date };
		}
		return { label: String(date), value: String(date) };
	});
}

const MONTHS = [
	'Januari',
	'Februari',
	'Maret',
	'April',
	'Mei',
	'Juni',
	'Juli',
	'Agustus',
	'September',
	'Oktober',
	'November',
	'Desember',
];

export function createOptionsOfMonths() {
	return MONTHS.map((month, idx) => {
		if (idx < 9) {
			// if idx is less than 10, add 0 before the idx
			return { label: month, value: '0' + (idx + 1) };
		}
		return { label: month, value: String(idx + 1) };
	});
}

function generateRangeOfYears() {
	const currentYear = new Date().getFullYear();
	const years = [];
	let startYear = 1980;
	while (startYear <= currentYear) {
		years.push(startYear);
		startYear++;
	}
	return years;
}

const YEARS = generateRangeOfYears();

export function createOptionOfYears() {
	return YEARS.map(year => ({ label: String(year), value: String(year) }));
}
