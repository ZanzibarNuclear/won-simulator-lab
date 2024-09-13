function formatDateTime(date) {
  const hours = String(date.getHours()).padStart(2, "0");
  const minutes = String(date.getMinutes()).padStart(2, "0");
  const timezone = date
    .toLocaleTimeString("en-us", { timeZoneName: "short" })
    .split(" ")[2];
  const day = String(date.getDate()).padStart(2, "0");
  const months = [
    "Jan",
    "Feb",
    "Mar",
    "Apr",
    "May",
    "Jun",
    "Jul",
    "Aug",
    "Sep",
    "Oct",
    "Nov",
    "Dec",
  ];
  const month = months[date.getMonth()];
  const year = date.getFullYear();

  return `${year}-${month}-${day} at ${hours}:${minutes} ${timezone}`;
  //   return `${hours}:${minutes} on ${day}-${month}-${year}`;
}
