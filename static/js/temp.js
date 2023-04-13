const submit_button = document.getElementById("reservation-availability");
let attention = Prompt();
let htmlContent = `
				<form
						action=""
						method="post"
						class="needs-validation"
						novalidate
						id="check-availability-form"
					>
						<div class="row" >
							<div class="col">
								<div class="row" id="reservation-dates-modal">
									<div class="col">
										<input
										disabled

											id="start"
											required
											type="text"
											class="form-control"
											aria-describedby="start date"
											placeholder="Arrival date"
										/>
									</div>
									<div class="col">
										<input
										disabled
											required
											type="text"
											class="form-control"
											id="end"
											aria-describedby="end date"
											placeholder="Check-out date"
										/>
									</div>
								</div>
							</div>
						</div>
					</form>`;

submit_button.addEventListener("click", () => {
     attention.customeHTML({
          msg: "Choose you dates",
          htmlContent: htmlContent,
          callback: () => {
               let form = document.getElementById("check-availability-form");
               let csrf_token = document.getElementById("csrf_token").value;
               let sd = document.getElementById("start").value;
               let ed = document.getElementById("end").value;
               let formData = new FormData(form);
               formData.append("csrf_token", csrf_token);
               formData.append("room_id", "1");
               formData.append("start", sd);
               formData.append("end", ed);
               console.log(formData);
               fetch("/search-availability-json", {
                    method: "post",
                    body: formData,
               })
                    .then((resp) => resp.json())
                    .then((data) => {
                         if (data.ok) console.log(data);
                         else console.log("not available");
                    });
          },
          didOpenParam: () => {
               document.getElementById("start").removeAttribute("disabled");
               document.getElementById("end").removeAttribute("disabled");
          },
          willOpenParam: () => {
               const elem = document.getElementById("reservation-dates-modal");
               const rangepicker = new DateRangePicker(elem, {
                    format: "dd-mm-yyyy",
                    minDate: new Date(),
               });
          },
          preConfirmParam: () => {
               return [
                    document.getElementById("start").value,
                    document.getElementById("end").value,
               ];
          },
     });
});
