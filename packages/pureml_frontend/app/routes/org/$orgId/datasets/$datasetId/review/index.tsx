import { MetaFunction, redirect } from "@remix-run/node";
import { Form, useLoaderData, useNavigate } from "@remix-run/react";
import Breadcrumbs from "~/components/Breadcrumbs";
import ReviewTabbar from "~/components/ReviewTabbar";
import Tabbar from "~/components/Tabbar";
import AvatarIcon from "~/components/ui/Avatar";
import { fetchDatasetReview } from "~/routes/api/datasets.server";
import { commitSession, getSession } from "~/session";

export const meta: MetaFunction = () => ({
  charset: "utf-8",
  title: "Dataset Review | PureML",
  viewport: "width=device-width,initial-scale=1",
});

export async function loader({ params, request }: any) {
  const session = await getSession(request.headers.get("Cookie"));
  const allReview = await fetchDatasetReview(
    session.get("orgId"),
    params.datasetId,
    session.get("accessToken")
  );

  return { allReview, params };
}

export async function action({ request }: any) {
  const formData = await request.formData();
  let option = Object.fromEntries(formData);
  const session = await getSession(request.headers.get("Cookie"));
  session.set("datasetName", option.datasetName);
  session.set("version", option.version);
  session.set("fromBranch", option.fromBranch);
  session.set("toBranch", option.toBranch);
  return redirect(`${option.reviewId}/datalineage`, {
    headers: {
      "Set-Cookie": await commitSession(session),
    },
  });
}

export default function DatasetReview() {
  const reviewData = useLoaderData();
  const navigate = useNavigate();
  return (
    <div id="datasetsReview">
      <div className="flex justify-center sticky top-0 bg-slate-50 w-full border-b border-slate-200">
        <div className="flex justify-between px-12 2xl:pr-0 w-full max-w-screen-2xl">
          <Breadcrumbs />
          <Tabbar intent="primaryDatasetTab" tab="review" fullWidth={false} />
        </div>
      </div>
      <div className="flex justify-center w-full">
        <div className="bg-slate-50 flex flex-col h-screen overflow-hidden w-full 2xl:max-w-screen-2xl">
          <ReviewTabbar intent="datasetReviewTab" tab="newcommits" />
          <div className="px-12 pt-8 w-full h-[75%] overflow-auto">
            <div className="w-2/3">
              {reviewData && (
                <>
                  {reviewData.allReview ? (
                    <div className="">
                      {reviewData.allReview.map(
                        (review: any, index: number) => (
                          <>
                            {!review.is_complete && !review.is_accepted && (
                              <Form method="post">
                                <input
                                  hidden
                                  readOnly
                                  name="reviewId"
                                  value={review.uuid}
                                />
                                <input
                                  hidden
                                  readOnly
                                  name="version"
                                  value={review.from_branch_version.version}
                                />
                                <input
                                  hidden
                                  readOnly
                                  name="fromBranch"
                                  value={review.from_branch.name}
                                />
                                <input
                                  hidden
                                  readOnly
                                  name="toBranch"
                                  value={review.to_branch.name}
                                />
                                <input
                                  hidden
                                  readOnly
                                  name="datasetName"
                                  value={review.dataset.name}
                                />
                                <button
                                  className="pb-6 w-full"
                                  key={index}
                                  onClick={() => {
                                    navigate(`${review.uuid}/datalineage`);
                                  }}
                                >
                                  <div className="hover:bg-slate-100 rounded-2xl flex justify-between p-4">
                                    <div className="flex items-center">
                                      <AvatarIcon>
                                        {review.created_by.name
                                          .charAt(0)
                                          .toUpperCase()}
                                      </AvatarIcon>
                                      <div className="text-sm text-slate-600 px-4">
                                        <a
                                          href={`/${review.created_by.handle}`}
                                          className="font-medium text-slate-800"
                                        >
                                          {review.created_by.handle}
                                        </a>{" "}
                                        submitted{" "}
                                        {review.from_branch_version.version} of{" "}
                                        <a
                                          href={`/org/${reviewData.params.orgId}/datasets/${review.dataset.name}`}
                                          className="font-medium text-slate-800"
                                        >
                                          {review.dataset.name}
                                        </a>{" "}
                                        from {review.from_branch.name} to{" "}
                                        {review.to_branch.name}
                                      </div>
                                    </div>
                                    {/* <div>Time</div> */}
                                  </div>
                                </button>
                              </Form>
                            )}
                          </>
                        )
                      )}
                    </div>
                  ) : (
                    <div>No reviews added yet</div>
                  )}
                </>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

// ############################ error boundary ###########################

export function ErrorBoundary() {
  return (
    <div className="flex flex-col h-screen justify-center items-center bg-slate-50">
      <div className="text-3xl text-slate-600 font-medium">Oops!!</div>
      <div className="text-3xl text-slate-600 font-medium">
        Something went wrong :(
      </div>
      <img src="/error/ErrorFunction.gif" alt="Error" width="500" />
    </div>
  );
}

export function CatchBoundary() {
  return (
    <div className="flex flex-col h-screen justify-center items-center bg-slate-50">
      <div className="text-3xl text-slate-600 font-medium">Oops!!</div>
      <div className="text-3xl text-slate-600 font-medium">
        Something went wrong :(
      </div>
      <img src="/error/ErrorFunction.gif" alt="Error" width="500" />
    </div>
  );
}
